package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/in"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

var planRowsRe = regexp.MustCompile(`rows=(\d+)`)

// WriteService implements in.WriteUseCase: preview → approve → execute.
type WriteService struct {
	gateway        *Gateway
	validator      out.Validator
	policy         out.PolicyEngine
	pools          map[string]out.Pool
	tokens         out.TokenStore
	approval       out.ApprovalEngine
	audit          out.AuditSink
	clock          out.Clock
	redact         out.Redactor
	tokenTTL       time.Duration
	queryTimeout   time.Duration
	maxQueryLength int
}

var _ in.WriteUseCase = (*WriteService)(nil)

// NewWriteService wires a WriteService.
func NewWriteService(
	gateway *Gateway,
	validator out.Validator,
	policy out.PolicyEngine,
	pools map[string]out.Pool,
	tokens out.TokenStore,
	approval out.ApprovalEngine,
	audit out.AuditSink,
	clock out.Clock,
	redact out.Redactor,
	tokenTTL time.Duration,
	queryTimeout time.Duration,
	maxQueryLength int,
) *WriteService {
	if pools == nil {
		pools = map[string]out.Pool{}
	}
	return &WriteService{
		gateway: gateway, validator: validator, policy: policy,
		pools: pools, tokens: tokens, approval: approval,
		audit: audit, clock: clock, redact: redact,
		tokenTTL: tokenTTL, queryTimeout: queryTimeout, maxQueryLength: maxQueryLength,
	}
}

// Preview estimates a write and issues a single-use token. Never executes.
func (s *WriteService) Preview(ctx context.Context, req domain.ReadRequest) (*in.WritePreview, error) {
	if err := s.gateway.RequireScope(ctx, domain.ScopeWritePreview); err != nil {
		return nil, err
	}
	if s.maxQueryLength > 0 && len(req.SQL) > s.maxQueryLength {
		return nil, domain.NewError(domain.CodeQueryTooLong, "query exceeds max_query_length")
	}
	vq, err := s.validator.Validate(ctx, req.Connection, req.SQL)
	if err != nil {
		s.emitDenied(ctx, req, "db_write_preview", err)
		return nil, err
	}
	decision := s.policy.Evaluate(ctx, vq)
	if decision.Action == domain.ActionDeny {
		err := domain.ErrQueryDenied(decision.Reason, decision.Rule)
		s.emitDenied(ctx, req, "db_write_preview", err)
		return nil, err
	}
	if !vq.StatementType.IsWrite() {
		err := domain.ErrQueryDenied("db_write_preview requires a write statement", decision.Rule)
		s.emitDenied(ctx, req, "db_write_preview", err)
		return nil, err
	}
	pool, ok := s.pools[req.Connection]
	if !ok {
		return nil, domain.ErrConnectionNotFound(req.Connection)
	}
	plan, err := pool.Explain(ctx, domain.Query{Connection: req.Connection, SQL: vq.NormalizedSQL, Params: req.Params})
	if err != nil {
		return nil, err
	}
	estimate := parseEstimate(plan.Text)
	warnings := buildWarnings(vq, estimate)

	id, err := randID("pvw_")
	if err != nil {
		return nil, domain.NewError(domain.CodeInternal, "generate preview token")
	}
	expiresAt := s.clock.Now()
	if s.tokenTTL > 0 {
		expiresAt = expiresAt.Add(s.tokenTTL)
	}
	token := &domain.PreviewToken{
		ID:          id,
		Connection:  req.Connection,
		SQLHash:     domain.HashSQL(vq.NormalizedSQL),
		ParamsHash:  domain.HashParams(req.Params),
		State:       domain.TokenPending,
		SQL:         vq.NormalizedSQL,
		Params:      req.Params,
		PlanText:    plan.Text,
		Warnings:    warnings,
		Statement:   vq.StatementType,
		RequestedBy: domain.TokenIDFrom(ctx),
		Estimate:    estimate,
	}
	if err := s.tokens.Create(ctx, token); err != nil {
		return nil, err
	}
	if s.audit != nil {
		s.audit.Emit(ctx, domain.AuditEvent{
			TS: s.clock.Now(), Event: "write.preview_created",
			RequestID: domain.RequestIDFrom(ctx), TokenID: domain.TokenIDFrom(ctx),
			Connection: req.Connection, Tool: "db_write_preview",
			SQLHash: token.SQLHash, SQLNormalized: s.redact.RedactString(vq.NormalizedSQL),
			StatementType: vq.StatementType, Tables: vq.Tables,
			PolicyDecision: string(decision.Action), PolicyRule: decision.Rule,
			RowCount: int(estimate), Status: "success",
		})
	}
	return &in.WritePreview{
		PreviewTokenID: id, ExpiresAt: expiresAt,
		StatementType: vq.StatementType, AffectedEstimate: estimate,
		Plan: plan.Text, Warnings: warnings, SQLNormalized: vq.NormalizedSQL,
	}, nil
}

// Execute verifies the single-use token and runs the write.
func (s *WriteService) Execute(ctx context.Context, req domain.ExecuteRequest) (*in.WriteResult, error) {
	if err := s.gateway.RequireScope(ctx, domain.ScopeWriteExecute); err != nil {
		return nil, err
	}
	token, err := s.tokens.Get(ctx, req.PreviewTokenID)
	if err != nil {
		return nil, err
	}
	switch token.State {
	case domain.TokenUsed:
		return nil, domain.NewError(domain.CodeTokenAlreadyUsed, "preview token already used")
	case domain.TokenExpired:
		return nil, domain.NewError(domain.CodeTokenExpired, "preview token expired")
	case domain.TokenRejected:
		return nil, domain.NewError(domain.CodeTokenMismatch, "write request was rejected: "+token.DecidedReason)
	}
	if token.State == domain.TokenPending {
		approved, err := s.approval.Approve(ctx, out.ApprovalRequest{
			Connection: token.Connection, Statement: token.Statement,
			SQL: token.SQL, Params: token.Params, Estimate: 0,
		})
		if err != nil {
			return nil, err
		}
		if !approved {
			_, _ = s.tokens.Reject(ctx, token.ID, domain.TokenIDFrom(ctx), "rejected by approver")
			return nil, domain.NewError(domain.CodeTokenMismatch, "write request was rejected by approver")
		}
		if _, err := s.tokens.Approve(ctx, token.ID, domain.TokenIDFrom(ctx), "approved"); err != nil {
			return nil, err
		}
	}
	// Re-validate: the SQL must byte-match the previewed statement.
	vq, err := s.validator.Validate(ctx, req.Connection, req.SQL)
	if err != nil {
		return nil, err
	}
	if req.Connection != token.Connection ||
		domain.HashSQL(vq.NormalizedSQL) != token.SQLHash ||
		domain.HashParams(req.Params) != token.ParamsHash {
		return nil, domain.NewError(domain.CodeTokenMismatch, "preview token does not match connection, SQL or params")
	}
	decision := s.policy.Evaluate(ctx, vq)
	if decision.Action == domain.ActionDeny {
		return nil, domain.ErrQueryDenied(decision.Reason, decision.Rule)
	}
	if err := s.tokens.MarkUsed(ctx, token.ID); err != nil {
		return nil, err
	}
	pool, ok := s.pools[req.Connection]
	if !ok {
		return nil, domain.ErrConnectionNotFound(req.Connection)
	}
	qctx := ctx
	if s.queryTimeout > 0 {
		var cancel context.CancelFunc
		qctx, cancel = context.WithTimeout(ctx, s.queryTimeout)
		defer cancel()
	}
	start := s.clock.Now()
	res, err := pool.Exec(qctx, domain.Query{Connection: req.Connection, SQL: vq.NormalizedSQL, Params: req.Params})
	duration := s.clock.Now().Sub(start).Milliseconds()
	if err != nil {
		if qctx.Err() == context.DeadlineExceeded {
			return nil, domain.NewError(domain.CodeQueryTimeout, "query timed out")
		}
		return nil, err
	}
	if s.audit != nil {
		s.audit.Emit(ctx, domain.AuditEvent{
			TS: s.clock.Now(), Event: "write.executed",
			RequestID: domain.RequestIDFrom(ctx), TokenID: domain.TokenIDFrom(ctx),
			Connection: req.Connection, Tool: "db_write_execute",
			SQLHash: token.SQLHash, SQLNormalized: s.redact.RedactString(vq.NormalizedSQL),
			StatementType: vq.StatementType, Tables: vq.Tables,
			PolicyDecision: string(decision.Action), PolicyRule: decision.Rule,
			RowCount: int(res.RowsAffected), DurationMs: duration, Status: "success",
		})
	}
	return &in.WriteResult{RowsAffected: res.RowsAffected, DurationMs: duration}, nil
}

func (s *WriteService) emitDenied(ctx context.Context, req domain.ReadRequest, tool string, err error) {
	if s.audit == nil {
		return
	}
	s.audit.Emit(ctx, domain.AuditEvent{
		TS: s.clock.Now(), Event: "query.denied",
		RequestID: domain.RequestIDFrom(ctx), TokenID: domain.TokenIDFrom(ctx),
		Connection: req.Connection, Tool: tool,
		SQLHash: domain.HashSQL(req.SQL), SQLNormalized: s.redact.RedactString(req.SQL),
		Status: "denied", Error: err.Error(),
	})
}

// parseEstimate reads the largest rows=N from an EXPLAIN plan.
// The top node of a write plan reports rows=0, so the maximum across all
// nodes is the useful estimate.
func parseEstimate(plan string) int64 {
	var best int64
	for _, m := range planRowsRe.FindAllStringSubmatch(plan, -1) {
		var n int64
		_, _ = fmt.Sscanf(m[1], "%d", &n)
		if n > best {
			best = n
		}
	}
	return best
}

func buildWarnings(vq *domain.ValidatedQuery, estimate int64) []string {
	table := "?"
	if len(vq.Tables) > 0 {
		table = vq.Tables[0]
	}
	return []string{fmt.Sprintf("This will affect ~%d rows in table '%s'", estimate, table)}
}

func randID(prefix string) (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(b[:]), nil
}
