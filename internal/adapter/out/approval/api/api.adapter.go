package api

import (
	"context"
	"fmt"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// pollInterval is the fallback re-check cadence. Watch broadcasts are
// non-blocking and may drop under load, so polling guarantees progress.
const pollInterval = 2 * time.Second

// Engine waits for a human decision made via the Admin API/UI instead of
// prompting on a terminal. It never touches stdin, so it is safe for
// headless and container operation (approval.method: api).
//
// The wait is bounded by the earliest of the preview token expiry and the
// caller context deadline. If neither exists the call fails fast with an
// approval-required error instead of hanging forever.
type Engine struct {
	tokens out.TokenStore
	clock  out.Clock
}

var _ out.ApprovalEngine = (*Engine)(nil)

// NewApprovalEngine builds an API approval engine over the preview token store.
func NewApprovalEngine(tokens out.TokenStore, clock out.Clock) *Engine {
	return &Engine{tokens: tokens, clock: clock}
}

// Approve blocks until the pending preview token is approved, rejected,
// expired, or the wait deadline hits.
func (e *Engine) Approve(ctx context.Context, req out.ApprovalRequest) (bool, error) {
	if req.PreviewTokenID == "" {
		return false, fmt.Errorf("api approval requires a preview token id")
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	ch, cancel := e.tokens.Watch()
	defer cancel()
	for {
		t, err := e.tokens.Get(ctx, req.PreviewTokenID)
		if err != nil {
			return false, err
		}
		switch t.State {
		case domain.TokenApproved:
			return true, nil
		case domain.TokenRejected:
			return false, nil
		case domain.TokenUsed:
			return false, domain.NewError(domain.CodeTokenAlreadyUsed, "preview token already used")
		case domain.TokenExpired:
			return false, domain.NewError(domain.CodeTokenExpired, "preview token expired")
		}
		wait, ok := e.waitBound(ctx, t)
		if !ok {
			return false, domain.NewError(domain.CodeApprovalRequired,
				"write approval is pending; approve it via the Admin API and retry")
		}
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return false, ctx.Err()
		case ev := <-ch:
			timer.Stop()
			if ev.Token == nil || ev.Token.ID != req.PreviewTokenID {
				continue
			}
		case <-timer.C:
		}
	}
}

// waitBound returns how long to wait before re-checking, and whether waiting
// is bounded at all. Preference: token expiry, then context deadline.
func (e *Engine) waitBound(ctx context.Context, t *domain.PreviewToken) (time.Duration, bool) {
	wait := pollInterval
	bounded := false
	if !t.ExpiresAt.IsZero() {
		if remaining := t.ExpiresAt.Sub(e.clock.Now()); remaining > 0 {
			wait = min(wait, remaining)
			bounded = true
		} else {
			// Expiry is due; loop once more so Get surfaces it.
			return 0, true
		}
	}
	if deadline, ok := ctx.Deadline(); ok {
		if remaining := time.Until(deadline); remaining > 0 {
			wait = min(wait, remaining)
			bounded = true
		} else {
			return 0, true
		}
	}
	return wait, bounded
}
