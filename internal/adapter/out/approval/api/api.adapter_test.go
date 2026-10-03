package api

import (
	"context"
	"testing"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/adapter/out/token/memory"
	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

type mutableClock struct{ now time.Time }

func (c *mutableClock) Now() time.Time { return c.now }

func setup() (*Engine, *memory.Store, *mutableClock) {
	clk := &mutableClock{now: time.Now()}
	store := memory.New(clk, time.Minute)
	return NewApprovalEngine(store, clk), store, clk
}

func pendingToken(t *testing.T, store *memory.Store, clk *mutableClock) string {
	t.Helper()
	id := "pvw_test123"
	if err := store.Create(context.Background(), &domain.PreviewToken{
		ID: id, Connection: "analytics", SQL: "UPDATE orders SET a = 1",
		CreatedAt: clk.now,
	}); err != nil {
		t.Fatalf("seed pending token: %v", err)
	}
	return id
}

func TestApproveAlreadyApproved(t *testing.T) {
	eng, store, clk := setup()
	id := pendingToken(t, store, clk)
	if _, err := store.Approve(context.Background(), id, "boss", "ok"); err != nil {
		t.Fatalf("approve: %v", err)
	}
	ok, err := eng.Approve(context.Background(), out.ApprovalRequest{PreviewTokenID: id})
	if err != nil || !ok {
		t.Fatalf("Approve = %v, %v; want true, nil", ok, err)
	}
}

func TestApproveWaitsForAPIDecision(t *testing.T) {
	eng, store, clk := setup()
	id := pendingToken(t, store, clk)
	done := make(chan bool, 1)
	go func() {
		ok, err := eng.Approve(context.Background(), out.ApprovalRequest{PreviewTokenID: id})
		if err != nil {
			t.Errorf("Approve: %v", err)
			done <- false
			return
		}
		done <- ok
	}()
	select {
	case <-done:
		t.Fatal("Approve returned before any decision")
	case <-time.After(50 * time.Millisecond):
	}
	if _, err := store.Approve(context.Background(), id, "boss", "looks good"); err != nil {
		t.Fatalf("api approve: %v", err)
	}
	select {
	case ok := <-done:
		if !ok {
			t.Fatal("Approve = false; want true after API approval")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Approve did not unblock after API approval")
	}
}

func TestApproveRejected(t *testing.T) {
	eng, store, clk := setup()
	id := pendingToken(t, store, clk)
	if _, err := store.Reject(context.Background(), id, "boss", "too broad"); err != nil {
		t.Fatalf("reject: %v", err)
	}
	ok, err := eng.Approve(context.Background(), out.ApprovalRequest{PreviewTokenID: id})
	if err != nil || ok {
		t.Fatalf("Approve = %v, %v; want false, nil", ok, err)
	}
}

func TestApproveExpired(t *testing.T) {
	eng, store, clk := setup()
	id := pendingToken(t, store, clk)
	clk.now = clk.now.Add(2 * time.Minute) // past the 1m token TTL
	if _, err := eng.Approve(context.Background(), out.ApprovalRequest{PreviewTokenID: id}); err == nil {
		t.Fatal("Approve on expired token: want error")
	}
}

func TestApproveNoTokenID(t *testing.T) {
	eng, _, _ := setup()
	if _, err := eng.Approve(context.Background(), out.ApprovalRequest{}); err == nil {
		t.Fatal("Approve without token id: want error")
	}
}

func TestApproveUnboundedFailsFast(t *testing.T) {
	clk := &mutableClock{now: time.Now()}
	store := memory.New(clk, 0) // no expiry
	eng := NewApprovalEngine(store, clk)
	id := "pvw_noexpiry"
	if err := store.Create(context.Background(), &domain.PreviewToken{
		ID: id, Connection: "analytics", SQL: "UPDATE orders SET a = 1",
		CreatedAt: clk.now,
	}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	start := time.Now()
	_, err := eng.Approve(context.Background(), out.ApprovalRequest{PreviewTokenID: id})
	if err == nil || time.Since(start) > 2*time.Second {
		t.Fatalf("unbounded wait: want fast approval-required error, got %v", err)
	}
}

func TestApproveContextCancel(t *testing.T) {
	eng, store, clk := setup()
	id := pendingToken(t, store, clk)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := eng.Approve(ctx, out.ApprovalRequest{PreviewTokenID: id})
		done <- err
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("want context error after cancel")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Approve did not respect context cancellation")
	}
}
