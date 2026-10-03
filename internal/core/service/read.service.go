package service

import (
	"context"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/in"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// ReadService implements in.ReadUseCase: validate → policy → execute → audit.
type ReadService struct {
	gateway        *Gateway
	validator      out.Validator
	policy         out.PolicyEngine
	pools          map[string]out.Pool
	metas          map[string]domain.ConnectionMeta
	audit          out.AuditSink
	clock          out.Clock
	redact         out.Redactor
	rowLimit       int
	maxQueryLength int
	queryTimeout   time.Duration
}

var _ in.ReadUseCase = (*ReadService)(nil)

// NewReadService wires a ReadService. Only *Gateway, validator, policy and
// audit are required; pools/metas may be empty (queries then fail with
// CONNECTION_NOT_FOUND until connectors are wired).
func NewReadService(
	gateway *Gateway,
	validator out.Validator,
	policy out.PolicyEngine,
	pools map[string]out.Pool,
	metas map[string]domain.ConnectionMeta,
	audit out.AuditSink,
	clock out.Clock,
	redact out.Redactor,
	rowLimit int,
	maxQueryLength int,
	queryTimeout time.Duration,
) *ReadService {
	if pools == nil {
		pools = map[string]out.Pool{}
	}
	if metas == nil {
		metas = map[string]domain.ConnectionMeta{}
	}
	return &ReadService{
		gateway: gateway, validator: validator, policy: policy,
		pools: pools, metas: metas, audit: audit, clock: clock, redact: redact,
		rowLimit: rowLimit, maxQueryLength: maxQueryLength, queryTimeout: queryTimeout,
	}
}

// Read executes a read-only query through the full pipeline.
func (s *ReadService) Read(ctx context.Context, req domain.ReadRequest) (*domain.ResultSet, error) {
	if err := s.gateway.RequireScope(ctx, domain.ScopeRead); err != nil {
		return nil, err
	}
	if s.maxQueryLength > 0 && len(req.SQL) > s.maxQueryLength {
		return nil, domain.NewError(domain.CodeQueryTooLong, "query exceeds max_query_length")
	}
	vq, err := s.validator.Validate(ctx, req.Connection, req.SQL)
	if err != nil {
		s.emitDenied(ctx, req, nil, err)
		return nil, err
	}
	decision := s.policy.Evaluate(ctx, vq)
	if decision.Action == domain.ActionDeny {
		err := domain.ErrQueryDenied(decision.Reason, decision.Rule)
		s.emitDenied(ctx, req, vq, err)
		return nil, err
	}
	if !vq.StatementType.IsRead() {
		if decision.Action == domain.ActionRequireApproval {
			return nil, domain.NewError(domain.CodeApprovalRequired, "write statements require db_write_preview first")
		}
		err := domain.ErrQueryDenied("statement type not allowed for db_read", decision.Rule)
		s.emitDenied(ctx, req, vq, err)
		return nil, err
	}
	pool, ok := s.pools[req.Connection]
	if !ok {
		return nil, domain.ErrConnectionNotFound(req.Connection)
	}
	rowLimit := s.rowLimitFor(req.Connection)
	limit := req.Limit
	if limit <= 0 || limit > rowLimit {
		limit = rowLimit
	}
	qctx := ctx
	if s.queryTimeout > 0 {
		var cancel context.CancelFunc
		qctx, cancel = context.WithTimeout(ctx, s.queryTimeout)
		defer cancel()
	}
	start := s.clock.Now()
	res, err := pool.Query(qctx, domain.Query{
		Connection: req.Connection,
		SQL:        vq.NormalizedSQL,
		Params:     req.Params,
		Limit:      limit,
	}, domain.QueryOpts{ReadOnly: true, Timeout: s.queryTimeout, RowLimit: rowLimit})
	duration := s.clock.Now().Sub(start).Milliseconds()
	if err != nil {
		if qctx.Err() == context.DeadlineExceeded {
			return nil, domain.NewError(domain.CodeQueryTimeout, "query timed out")
		}
		s.emitExecuted(ctx, req, vq, decision, 0, false, duration, "error", err.Error())
		return nil, err
	}
	s.emitExecuted(ctx, req, vq, decision, res.RowCount, res.Truncated, duration, "success", "")
	return res, nil
}

// Explain runs EXPLAIN without executing the query.
func (s *ReadService) Explain(ctx context.Context, req domain.ReadRequest) (*domain.Plan, error) {
	if err := s.gateway.RequireScope(ctx, domain.ScopeRead); err != nil {
		return nil, err
	}
	vq, err := s.validator.Validate(ctx, req.Connection, req.SQL)
	if err != nil {
		s.emitDenied(ctx, req, nil, err)
		return nil, err
	}
	decision := s.policy.Evaluate(ctx, vq)
	if decision.Action == domain.ActionDeny {
		err := domain.ErrQueryDenied(decision.Reason, decision.Rule)
		s.emitDenied(ctx, req, vq, err)
		return nil, err
	}
	pool, ok := s.pools[req.Connection]
	if !ok {
		return nil, domain.ErrConnectionNotFound(req.Connection)
	}
	return pool.Explain(ctx, domain.Query{Connection: req.Connection, SQL: vq.NormalizedSQL, Params: req.Params})
}

func (s *ReadService) rowLimitFor(connection string) int {
	if m, ok := s.metas[connection]; ok && m.RowLimit > 0 {
		return m.RowLimit
	}
	if s.rowLimit > 0 {
		return s.rowLimit
	}
	return 1000
}

func (s *ReadService) emitDenied(ctx context.Context, req domain.ReadRequest, vq *domain.ValidatedQuery, err error) {
	if s.audit == nil {
		return
	}
	e := domain.AuditEvent{
		TS:         s.clock.Now(),
		Event:      "query.denied",
		RequestID:  domain.RequestIDFrom(ctx),
		TokenID:    domain.TokenIDFrom(ctx),
		Connection: req.Connection,
		Tool:       "db_read",
		Status:     "denied",
		Error:      err.Error(),
	}
	if vq != nil {
		e.SQLHash = domain.HashSQL(vq.NormalizedSQL)
		e.SQLNormalized = s.redact.RedactString(vq.NormalizedSQL)
		e.StatementType = vq.StatementType
		e.Tables = vq.Tables
	} else {
		e.SQLHash = domain.HashSQL(req.SQL)
		e.SQLNormalized = s.redact.RedactString(req.SQL)
	}
	if de, ok := err.(*domain.Error); ok {
		e.PolicyRule = de.Rule
	}
	s.audit.Emit(ctx, e)
}

func (s *ReadService) emitExecuted(ctx context.Context, req domain.ReadRequest, vq *domain.ValidatedQuery, decision domain.PolicyDecision, rows int, truncated bool, duration int64, status, errMsg string) {
	if s.audit == nil {
		return
	}
	s.audit.Emit(ctx, domain.AuditEvent{
		TS:             s.clock.Now(),
		Event:          "query.executed",
		RequestID:      domain.RequestIDFrom(ctx),
		TokenID:        domain.TokenIDFrom(ctx),
		Connection:     req.Connection,
		Tool:           "db_read",
		SQLHash:        domain.HashSQL(vq.NormalizedSQL),
		SQLNormalized:  s.redact.RedactString(vq.NormalizedSQL),
		StatementType:  vq.StatementType,
		Tables:         vq.Tables,
		PolicyDecision: string(decision.Action),
		PolicyRule:     decision.Rule,
		RowCount:       rows,
		Truncated:      truncated,
		DurationMs:     duration,
		Status:         status,
		Error:          errMsg,
	})
}
