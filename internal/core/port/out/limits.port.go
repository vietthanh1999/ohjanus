package out

import "context"

// RateLimiter enforces per-key request rates (token bucket).
// Implemented by adapter/out/limits/*.
type RateLimiter interface {
	// Allow reports whether a request for key (usually the token id)
	// may proceed, consuming one token when true.
	Allow(ctx context.Context, key string) bool
}

// Metrics records Prometheus metrics (§13.1 of the spec).
// Implemented by adapter/out/metrics/*.
type Metrics interface {
	IncRequests(tool, status string)
	ObserveQueryDuration(connection, statement string, seconds float64)
	IncPolicyDenials(rule, connection string)
}

// NoopMetrics discards all metrics (default when disabled).
type NoopMetrics struct{}

// IncRequests implements Metrics.
func (NoopMetrics) IncRequests(string, string) {}

// ObserveQueryDuration implements Metrics.
func (NoopMetrics) ObserveQueryDuration(string, string, float64) {}

// IncPolicyDenials implements Metrics.
func (NoopMetrics) IncPolicyDenials(string, string) {}

// AllowRateLimiter always allows (rate limiting disabled).
type AllowRateLimiter struct{}

// Allow implements RateLimiter.
func (AllowRateLimiter) Allow(context.Context, string) bool { return true }
