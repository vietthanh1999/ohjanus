package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

type denyLimiter struct{}

func (denyLimiter) Allow(context.Context, string) bool { return false }

type countMetrics struct {
	requests map[string]int
	denials  int
}

func (c *countMetrics) IncRequests(tool, status string) {
	if c.requests == nil {
		c.requests = map[string]int{}
	}
	c.requests[tool+"/"+status]++
}
func (c *countMetrics) ObserveQueryDuration(string, string, float64) {}
func (c *countMetrics) IncPolicyDenials(string, string)              { c.denials++ }

func TestEnterRateLimited(t *testing.T) {
	g := NewGatewayWithLimits(fixedClock{}, denyLimiter{}, 0, out.NoopMetrics{})
	_, err := g.Enter(domain.WithAuth(context.Background(), "t", []domain.Scope{domain.ScopeRead}))
	var de *domain.Error
	if !errors.As(err, &de) || de.Code != domain.CodeRateLimited {
		t.Errorf("err = %v, want RATE_LIMITED", err)
	}
}

func TestEnterConcurrency(t *testing.T) {
	g := NewGatewayWithLimits(fixedClock{}, out.AllowRateLimiter{}, 1, out.NoopMetrics{})
	ctx := domain.WithAuth(context.Background(), "t", []domain.Scope{domain.ScopeRead})
	release, err := g.Enter(ctx)
	if err != nil {
		t.Fatalf("first Enter: %v", err)
	}
	if _, err := g.Enter(ctx); err == nil {
		t.Error("second concurrent Enter should fail")
	} else {
		var de *domain.Error
		if !errors.As(err, &de) || de.Code != domain.CodeTooManyConcurrent {
			t.Errorf("err = %v, want TOO_MANY_CONCURRENT_QUERIES", err)
		}
	}
	release()
	if _, err := g.Enter(ctx); err != nil {
		t.Errorf("Enter after release: %v", err)
	}
}

func TestReadObservesMetrics(t *testing.T) {
	m := &countMetrics{}
	clk := fixedClock{t: time.Now()}
	vq := &domain.ValidatedQuery{StatementType: domain.StatementSelect, NormalizedSQL: "SELECT 1"}
	svc := NewReadService(
		NewGatewayWithLimits(clk, out.AllowRateLimiter{}, 0, m),
		&fakeValidator{vq: vq},
		&fakePolicy{decision: allowDecision()},
		NewConnRegistry(
			map[string]out.Pool{"analytics": &fakePool{res: &domain.ResultSet{RowCount: 1}}},
			map[string]domain.ConnectionMeta{"analytics": {}},
		),
		&captureAudit{}, clk, identRedact{}, 100, 10000, 0,
	)
	if _, err := svc.Read(authedCtx(), domain.ReadRequest{Connection: "analytics", SQL: "SELECT 1"}); err != nil {
		t.Fatal(err)
	}
	if m.requests["db_read/success"] != 1 {
		t.Errorf("requests = %v", m.requests)
	}

	svc2 := NewReadService(
		NewGatewayWithLimits(clk, out.AllowRateLimiter{}, 0, m),
		&fakeValidator{vq: &domain.ValidatedQuery{StatementType: domain.StatementDrop, NormalizedSQL: "DROP TABLE x"}},
		&fakePolicy{decision: domain.PolicyDecision{Action: domain.ActionDeny, Rule: "deny-ddl"}},
		NewConnRegistry(map[string]out.Pool{}, map[string]domain.ConnectionMeta{}),
		&captureAudit{}, clk, identRedact{}, 100, 10000, 0,
	)
	svc2.Read(authedCtx(), domain.ReadRequest{Connection: "analytics", SQL: "DROP TABLE x"})
	if m.requests["db_read/denied"] != 1 || m.denials != 1 {
		t.Errorf("requests = %v denials = %d", m.requests, m.denials)
	}
}
