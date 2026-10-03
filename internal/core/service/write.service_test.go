package service

import (
	"context"
	"errors"
	"testing"
	"time"

	tokememory "github.com/vietthanh1999/ohjanus/internal/adapter/out/token/memory"
	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/in"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

type writeValidator struct {
	vq *domain.ValidatedQuery
}

func (f *writeValidator) Validate(context.Context, string, string) (*domain.ValidatedQuery, error) {
	return f.vq, nil
}

type writePolicy struct{}

func (writePolicy) Evaluate(context.Context, *domain.ValidatedQuery) domain.PolicyDecision {
	return domain.PolicyDecision{Action: domain.ActionRequireApproval, Rule: "warn-write"}
}

type writePool struct {
	planText string
	affected int64
}

func (p *writePool) Query(context.Context, domain.Query, domain.QueryOpts) (*domain.ResultSet, error) {
	return nil, errors.New("no reads here")
}
func (p *writePool) Exec(context.Context, domain.Query) (*domain.ExecResult, error) {
	return &domain.ExecResult{RowsAffected: p.affected}, nil
}
func (p *writePool) Explain(context.Context, domain.Query) (*domain.Plan, error) {
	return &domain.Plan{Text: p.planText}, nil
}
func (p *writePool) Schema(context.Context, string, string) ([]domain.Schema, error) {
	return nil, nil
}
func (p *writePool) Ping(context.Context) error { return nil }
func (p *writePool) Close() error               { return nil }

type fakeApproval struct {
	approved bool
	calls    int
}

func (f *fakeApproval) Approve(context.Context, out.ApprovalRequest) (bool, error) {
	f.calls++
	return f.approved, nil
}

type mutableClock struct{ now time.Time }

func (c *mutableClock) Now() time.Time { return c.now }

func writeSetup(approved bool) (*WriteService, *captureAudit, *fakeApproval, *mutableClock) {
	clk := &mutableClock{now: time.Now()}
	audit := &captureAudit{}
	vq := &domain.ValidatedQuery{
		StatementType: domain.StatementUpdate,
		Tables:        []string{"orders"},
		NormalizedSQL: "UPDATE orders SET a = 1",
	}
	svc := NewWriteService(
		NewGateway(clk),
		&writeValidator{vq: vq},
		writePolicy{},
		map[string]out.Pool{"analytics": &writePool{planText: "Update on orders  (cost=0.00..1.42 rows=42 width=6)", affected: 42}},
		tokememory.New(clk, time.Minute),
		&fakeApproval{approved: approved},
		audit, clk, identRedact{}, time.Minute, 30*time.Second, 10000,
	)
	return svc, audit, svc.approval.(*fakeApproval), clk
}

func writeCtx(scopes ...domain.Scope) context.Context {
	return domain.WithAuth(context.Background(), "tok_w", scopes)
}

func TestWritePreviewExecute(t *testing.T) {
	svc, audit, approval, _ := writeSetup(true)
	ctx := writeCtx(domain.ScopeWritePreview, domain.ScopeWriteExecute)
	req := domain.ReadRequest{Connection: "analytics", SQL: "UPDATE orders SET a = 1"}

	prev, err := svc.Preview(ctx, req)
	if err != nil {
		t.Fatalf("Preview: %v", err)
	}
	if prev.AffectedEstimate != 42 {
		t.Errorf("estimate = %d, want 42", prev.AffectedEstimate)
	}
	if len(prev.Warnings) != 1 {
		t.Errorf("warnings = %v", prev.Warnings)
	}

	exec, err := svc.Execute(ctx, domain.ExecuteRequest{Connection: "analytics", SQL: "UPDATE orders SET a = 1", PreviewTokenID: prev.PreviewTokenID})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if exec.RowsAffected != 42 {
		t.Errorf("rows = %d, want 42", exec.RowsAffected)
	}
	if approval.calls != 1 {
		t.Errorf("approval calls = %d, want 1", approval.calls)
	}

	if _, err := svc.Execute(ctx, domain.ExecuteRequest{Connection: "analytics", SQL: "UPDATE orders SET a = 1", PreviewTokenID: prev.PreviewTokenID}); err == nil {
		t.Error("replay should fail")
	} else {
		var de *domain.Error
		if !errors.As(err, &de) || de.Code != domain.CodeTokenAlreadyUsed {
			t.Errorf("err = %v, want TOKEN_ALREADY_USED", err)
		}
	}

	var created, executed bool
	for _, e := range audit.events {
		if e.Event == "write.preview_created" {
			created = true
		}
		if e.Event == "write.executed" {
			executed = true
		}
	}
	if !created || !executed {
		t.Errorf("audit events missing preview/executed: %+v", audit.events)
	}
}

func TestWritePreviewRejectsReads(t *testing.T) {
	svc, _, _, _ := writeSetup(true)
	svc.validator = &writeValidator{vq: &domain.ValidatedQuery{StatementType: domain.StatementSelect, NormalizedSQL: "SELECT 1"}}
	_, err := svc.Preview(writeCtx(domain.ScopeWritePreview), domain.ReadRequest{Connection: "analytics", SQL: "SELECT 1"})
	var de *domain.Error
	if !errors.As(err, &de) || de.Code != domain.CodeQueryDenied {
		t.Errorf("err = %v, want QUERY_DENIED", err)
	}
}

func TestWriteExecuteMismatch(t *testing.T) {
	svc, _, _, _ := writeSetup(true)
	ctx := writeCtx(domain.ScopeWritePreview, domain.ScopeWriteExecute)
	prev, err := svc.Preview(ctx, domain.ReadRequest{Connection: "analytics", SQL: "UPDATE orders SET a = 1"})
	if err != nil {
		t.Fatal(err)
	}
	svc.validator = &writeValidator{vq: &domain.ValidatedQuery{StatementType: domain.StatementUpdate, NormalizedSQL: "UPDATE orders SET a = 2"}}
	_, err = svc.Execute(ctx, domain.ExecuteRequest{Connection: "analytics", SQL: "UPDATE orders SET a = 2", PreviewTokenID: prev.PreviewTokenID})
	var de *domain.Error
	if !errors.As(err, &de) || de.Code != domain.CodeTokenMismatch {
		t.Errorf("err = %v, want TOKEN_MISMATCH", err)
	}
}

func TestWriteExecuteRejected(t *testing.T) {
	svc, _, _, _ := writeSetup(false)
	ctx := writeCtx(domain.ScopeWritePreview, domain.ScopeWriteExecute)
	prev, err := svc.Preview(ctx, domain.ReadRequest{Connection: "analytics", SQL: "UPDATE orders SET a = 1"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.Execute(ctx, domain.ExecuteRequest{Connection: "analytics", SQL: "UPDATE orders SET a = 1", PreviewTokenID: prev.PreviewTokenID})
	var de *domain.Error
	if !errors.As(err, &de) || de.Code != domain.CodeTokenMismatch {
		t.Errorf("err = %v, want TOKEN_MISMATCH (rejected)", err)
	}
}

func TestWriteScopeEnforced(t *testing.T) {
	svc, _, _, _ := writeSetup(true)
	_, err := svc.Preview(writeCtx(domain.ScopeRead), domain.ReadRequest{Connection: "analytics", SQL: "UPDATE x SET a=1"})
	var de *domain.Error
	if !errors.As(err, &de) || de.Code != domain.CodeForbidden {
		t.Errorf("err = %v, want FORBIDDEN", err)
	}
	var _ = in.WritePreview{}
}
