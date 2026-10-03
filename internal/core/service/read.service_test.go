package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

type fakeValidator struct {
	vq  *domain.ValidatedQuery
	err error
}

func (f *fakeValidator) Validate(context.Context, string, string) (*domain.ValidatedQuery, error) {
	return f.vq, f.err
}

type fakePolicy struct {
	decision domain.PolicyDecision
}

func (f *fakePolicy) Evaluate(context.Context, *domain.ValidatedQuery) domain.PolicyDecision {
	return f.decision
}

type fakePool struct {
	res         *domain.ResultSet
	err         error
	gotReadOnly bool
}

func (f *fakePool) Query(_ context.Context, q domain.Query, opts domain.QueryOpts) (*domain.ResultSet, error) {
	f.gotReadOnly = opts.ReadOnly
	return f.res, f.err
}

func (f *fakePool) Explain(context.Context, domain.Query) (*domain.Plan, error) {
	return &domain.Plan{Text: "Seq Scan"}, nil
}

func (f *fakePool) Schema(context.Context, string, string) ([]domain.Schema, error) { return nil, nil }
func (f *fakePool) Ping(context.Context) error                                      { return nil }
func (f *fakePool) Close() error                                                    { return nil }

type captureAudit struct {
	events []domain.AuditEvent
}

func (c *captureAudit) Emit(_ context.Context, e domain.AuditEvent) { c.events = append(c.events, e) }
func (c *captureAudit) Close() error                                { return nil }

type fixedClock struct{ t time.Time }

func (c fixedClock) Now() time.Time { return c.t }

type identRedact struct{}

func (identRedact) RedactString(s string) string   { return s }
func (identRedact) RedactValue(_, v string) string { return v }

func testSetup(vq *domain.ValidatedQuery, verr error, decision domain.PolicyDecision, res *domain.ResultSet, perr error) (*ReadService, *captureAudit, *fakePool) {
	audit := &captureAudit{}
	pool := &fakePool{res: res, err: perr}
	clock := fixedClock{t: time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC)}
	svc := NewReadService(
		NewGateway(clock),
		&fakeValidator{vq: vq, err: verr},
		&fakePolicy{decision: decision},
		map[string]out.Pool{"analytics": pool},
		map[string]domain.ConnectionMeta{
			"analytics": {Connection: domain.Connection{Name: "analytics", Driver: "postgres", ReadOnly: true}, RowLimit: 100},
		},
		audit, clock, identRedact{}, 1000, 10000, 30*time.Second,
	)
	return svc, audit, pool
}

func authedCtx() context.Context {
	return domain.WithAuth(context.Background(), "tok_1", []domain.Scope{domain.ScopeRead})
}

func allowDecision() domain.PolicyDecision {
	return domain.PolicyDecision{Action: domain.ActionAllow, Rule: "allow-select"}
}

func TestReadSuccess(t *testing.T) {
	vq := &domain.ValidatedQuery{
		StatementType: domain.StatementSelect,
		Tables:        []string{"orders"},
		NormalizedSQL: "SELECT id FROM orders",
	}
	svc, audit, pool := testSetup(vq, nil, allowDecision(),
		&domain.ResultSet{Columns: []string{"id"}, Rows: [][]any{{1}}, RowCount: 1}, nil)
	res, err := svc.Read(authedCtx(), domain.ReadRequest{Connection: "analytics", SQL: "SELECT id FROM orders"})
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if res.RowCount != 1 {
		t.Errorf("row_count = %d, want 1", res.RowCount)
	}
	if !pool.gotReadOnly {
		t.Error("expected read-only execution")
	}
	if len(audit.events) != 1 || audit.events[0].Event != "query.executed" || audit.events[0].Status != "success" {
		t.Errorf("audit events = %+v, want one successful query.executed", audit.events)
	}
	if audit.events[0].TokenID != "tok_1" || len(audit.events[0].Tables) != 1 {
		t.Errorf("audit event = %+v", audit.events[0])
	}
}

func TestReadDenied(t *testing.T) {
	vq := &domain.ValidatedQuery{StatementType: domain.StatementDrop, NormalizedSQL: "DROP TABLE x"}
	svc, audit, _ := testSetup(vq, nil,
		domain.PolicyDecision{Action: domain.ActionDeny, Rule: "deny-ddl", Reason: "blocked"}, nil, nil)
	_, err := svc.Read(authedCtx(), domain.ReadRequest{Connection: "analytics", SQL: "DROP TABLE x"})
	var de *domain.Error
	if !errors.As(err, &de) || de.Code != domain.CodeQueryDenied {
		t.Fatalf("err = %v, want QUERY_DENIED", err)
	}
	if len(audit.events) != 1 || audit.events[0].Event != "query.denied" {
		t.Errorf("audit events = %+v, want one query.denied", audit.events)
	}
}

func TestReadUnauthenticated(t *testing.T) {
	svc, audit, _ := testSetup(nil, nil, allowDecision(), nil, nil)
	if _, err := svc.Read(context.Background(), domain.ReadRequest{Connection: "analytics", SQL: "SELECT 1"}); err == nil {
		t.Error("expected auth error, got nil")
	}
	if len(audit.events) != 0 {
		t.Errorf("expected no audit events, got %+v", audit.events)
	}
}

func TestReadWriteRequiresApproval(t *testing.T) {
	vq := &domain.ValidatedQuery{StatementType: domain.StatementUpdate, NormalizedSQL: "UPDATE orders SET a = 1"}
	svc, _, _ := testSetup(vq, nil,
		domain.PolicyDecision{Action: domain.ActionRequireApproval, Rule: "warn-write"}, nil, nil)
	_, err := svc.Read(authedCtx(), domain.ReadRequest{Connection: "analytics", SQL: "UPDATE orders SET a = 1"})
	var de *domain.Error
	if !errors.As(err, &de) || de.Code != domain.CodeApprovalRequired {
		t.Fatalf("err = %v, want APPROVAL_REQUIRED", err)
	}
}

func TestReadConnectionNotFound(t *testing.T) {
	vq := &domain.ValidatedQuery{StatementType: domain.StatementSelect, NormalizedSQL: "SELECT 1"}
	svc, _, _ := testSetup(vq, nil, allowDecision(), nil, nil)
	_, err := svc.Read(authedCtx(), domain.ReadRequest{Connection: "nope", SQL: "SELECT 1"})
	var de *domain.Error
	if !errors.As(err, &de) || de.Code != domain.CodeConnectionNotFound {
		t.Fatalf("err = %v, want CONNECTION_NOT_FOUND", err)
	}
}

func TestReadQueryTooLong(t *testing.T) {
	svc, _, _ := testSetup(nil, nil, allowDecision(), nil, nil)
	_, err := svc.Read(authedCtx(), domain.ReadRequest{Connection: "analytics", SQL: string(make([]byte, 10001))})
	var de *domain.Error
	if !errors.As(err, &de) || de.Code != domain.CodeQueryTooLong {
		t.Fatalf("err = %v, want QUERY_TOO_LONG", err)
	}
}
