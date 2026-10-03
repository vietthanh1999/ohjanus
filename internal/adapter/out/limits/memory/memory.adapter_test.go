package memory

import (
	"context"
	"testing"
	"time"
)

type manualClock struct{ now time.Time }

func (c *manualClock) Now() time.Time { return c.now }

func TestAllowBurstThenDeny(t *testing.T) {
	clk := &manualClock{now: time.Now()}
	l := New(60, clk) // 1/sec sustained, burst 60
	ctx := context.Background()
	for i := 0; i < 60; i++ {
		if !l.Allow(ctx, "tok") {
			t.Fatalf("request %d denied, want allow (burst)", i)
		}
	}
	if l.Allow(ctx, "tok") {
		t.Error("61st immediate request allowed, want deny")
	}
}

func TestRefill(t *testing.T) {
	clk := &manualClock{now: time.Now()}
	l := New(60, clk)
	ctx := context.Background()
	for i := 0; i < 60; i++ {
		l.Allow(ctx, "tok")
	}
	clk.now = clk.now.Add(30 * time.Second) // +30 tokens
	allowed := 0
	for i := 0; i < 31; i++ {
		if l.Allow(ctx, "tok") {
			allowed++
		}
	}
	if allowed != 30 {
		t.Errorf("allowed = %d, want 30 after 30s refill", allowed)
	}
}

func TestKeysAreIndependent(t *testing.T) {
	clk := &manualClock{now: time.Now()}
	l := New(1, clk)
	ctx := context.Background()
	l.Allow(ctx, "a")
	if l.Allow(ctx, "b") != true {
		t.Error("different key should have its own bucket")
	}
	if l.Allow(ctx, "a") {
		t.Error("same key exhausted should deny")
	}
}

func TestZeroDisables(t *testing.T) {
	clk := &manualClock{now: time.Now()}
	l := New(0, clk)
	for i := 0; i < 1000; i++ {
		if !l.Allow(context.Background(), "tok") {
			t.Fatal("rpm<=0 should always allow")
		}
	}
}
