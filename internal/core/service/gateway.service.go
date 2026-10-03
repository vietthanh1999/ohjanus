package service

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Gateway enforces cross-cutting concerns: scopes, rate limits,
// concurrency limits, and metrics.
// Auth identity travels on the context via domain.WithAuth (adapter/in sets
// it after verifying the token), so adapter/in never imports this package.
type Gateway struct {
	clock   out.Clock
	limiter out.RateLimiter
	sem     chan struct{}
	metrics out.Metrics
}

// NewGateway builds a Gateway with no rate, concurrency or metrics limits.
func NewGateway(clock out.Clock) *Gateway {
	return &Gateway{clock: clock, limiter: out.AllowRateLimiter{}, metrics: out.NoopMetrics{}}
}

// NewGatewayWithLimits builds a Gateway with rate limiting (requests per
// minute per token, <=0 disables), a concurrency semaphore (<=0 disables)
// and metrics (nil means discard).
func NewGatewayWithLimits(clock out.Clock, limiter out.RateLimiter, maxConcurrent int, metrics out.Metrics) *Gateway {
	if limiter == nil {
		limiter = out.AllowRateLimiter{}
	}
	if metrics == nil {
		metrics = out.NoopMetrics{}
	}
	g := &Gateway{clock: clock, limiter: limiter, metrics: metrics}
	if maxConcurrent > 0 {
		g.sem = make(chan struct{}, maxConcurrent)
	}
	return g
}

// RequireScope enforces token scopes; ScopeAdmin implies every scope.
func (g *Gateway) RequireScope(ctx context.Context, s domain.Scope) error {
	return domain.RequireScope(ctx, s)
}

// Enter enforces the rate limit and concurrency cap for one query.
// The returned release function must be called when the query finishes.
func (g *Gateway) Enter(ctx context.Context) (release func(), err error) {
	noop := func() {}
	if !g.limiter.Allow(ctx, domain.TokenIDFrom(ctx)) {
		return noop, domain.NewError(domain.CodeRateLimited, "rate limit exceeded")
	}
	if g.sem == nil {
		return noop, nil
	}
	select {
	case g.sem <- struct{}{}:
		return func() { <-g.sem }, nil
	default:
		return noop, domain.NewError(domain.CodeTooManyConcurrent, "too many concurrent queries")
	}
}

// Observe records a finished tool call.
func (g *Gateway) Observe(tool, status string) {
	g.metrics.IncRequests(tool, status)
}

// ObserveDenial records a policy denial.
func (g *Gateway) ObserveDenial(rule, connection string) {
	g.metrics.IncPolicyDenials(rule, connection)
}

// ObserveDuration records a query duration.
func (g *Gateway) ObserveDuration(connection, statement string, seconds float64) {
	g.metrics.ObserveQueryDuration(connection, statement, seconds)
}
