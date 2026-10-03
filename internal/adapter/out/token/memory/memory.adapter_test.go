package memory

import (
	"context"
	"testing"
	"time"

	systemclock "github.com/vietthanh1999/ohjanus/internal/adapter/out/clock/system"
	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

type manualClock struct{ now time.Time }

func (c *manualClock) Now() time.Time { return c.now }

func TestTokenLifecycle(t *testing.T) {
	clk := &manualClock{now: time.Now()}
	s := New(clk, time.Minute)
	ctx := context.Background()

	if err := s.Create(ctx, &domain.PreviewToken{ID: "pvw_1", Connection: "a", SQLHash: "h", ParamsHash: "p"}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := s.Create(ctx, &domain.PreviewToken{ID: "pvw_1"}); err == nil {
		t.Error("duplicate Create should fail")
	}

	got, err := s.Get(ctx, "pvw_1")
	if err != nil || got.State != domain.TokenPending {
		t.Fatalf("Get = %+v, %v", got, err)
	}

	if _, err := s.Reject(ctx, "pvw_1", "u", ""); err == nil {
		t.Error("Reject without reason should fail")
	}
	rej, err := s.Reject(ctx, "pvw_1", "u", "nope")
	if err != nil || rej.State != domain.TokenRejected {
		t.Fatalf("Reject = %+v, %v", rej, err)
	}
	if err := s.MarkUsed(ctx, "pvw_1"); err == nil {
		t.Error("MarkUsed on rejected should fail")
	}

	if err := s.Create(ctx, &domain.PreviewToken{ID: "pvw_2"}); err != nil {
		t.Fatalf("Create: %v", err)
	}
	appr, err := s.Approve(ctx, "pvw_2", "admin", "ok")
	if err != nil || appr.State != domain.TokenApproved || appr.ApprovedBy != "admin" {
		t.Fatalf("Approve = %+v, %v", appr, err)
	}
	// Idempotent re-approve.
	if _, err := s.Approve(ctx, "pvw_2", "admin", "ok"); err != nil {
		t.Fatalf("re-Approve: %v", err)
	}
	if err := s.MarkUsed(ctx, "pvw_2"); err != nil {
		t.Fatalf("MarkUsed: %v", err)
	}
	if err := s.MarkUsed(ctx, "pvw_2"); err == nil {
		t.Error("double MarkUsed should fail")
	}

	list, err := s.List(ctx)
	if err != nil || len(list) != 2 {
		t.Fatalf("List = %d, %v", len(list), err)
	}

	if _, err := s.Get(ctx, "missing"); err == nil {
		t.Error("Get unknown should fail")
	}
}

func TestTokenExpiry(t *testing.T) {
	clk := &manualClock{now: time.Now()}
	s := New(clk, time.Minute)
	ctx := context.Background()
	if err := s.Create(ctx, &domain.PreviewToken{ID: "pvw_x"}); err != nil {
		t.Fatal(err)
	}
	clk.now = clk.now.Add(2 * time.Minute)
	got, err := s.Get(ctx, "pvw_x")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.State != domain.TokenExpired {
		t.Errorf("state = %s, want expired", got.State)
	}
	if _, err := s.Approve(ctx, "pvw_x", "u", "r"); err == nil {
		t.Error("Approve expired should fail")
	}
}

func TestWatch(t *testing.T) {
	s := New(systemclock.Clock{}, time.Minute)
	ctx := context.Background()
	ch, cancel := s.Watch()
	defer cancel()
	if err := s.Create(ctx, &domain.PreviewToken{ID: "pvw_w"}); err != nil {
		t.Fatal(err)
	}
	select {
	case ev := <-ch:
		if ev.Type != domain.TokenEventCreated || ev.Token.ID != "pvw_w" {
			t.Errorf("event = %+v", ev)
		}
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for event")
	}
}
