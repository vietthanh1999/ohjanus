package memory

import (
	"context"
	"sync"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// bucket is one token-bucket: capacity burst, refilled rate tokens/sec.
type bucket struct {
	tokens float64
	last   time.Time
}

// Limiter is a per-key token-bucket rate limiter.
// requestsPerMinute sets both the sustained rate and the burst capacity.
type Limiter struct {
	mu                sync.Mutex
	buckets           map[string]*bucket
	requestsPerMinute int
	clock             out.Clock
}

var _ out.RateLimiter = (*Limiter)(nil)

// New builds a limiter. Non-positive rpm disables limiting (always allow).
func New(requestsPerMinute int, clock out.Clock) *Limiter {
	return &Limiter{
		buckets:           map[string]*bucket{},
		requestsPerMinute: requestsPerMinute,
		clock:             clock,
	}
}

// Allow consumes one token for key when available.
func (l *Limiter) Allow(_ context.Context, key string) bool {
	if l.requestsPerMinute <= 0 {
		return true
	}
	rate := float64(l.requestsPerMinute) / 60.0
	capacity := float64(l.requestsPerMinute)
	now := l.clock.Now()

	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: capacity, last: now}
		l.buckets[key] = b
	}
	elapsed := now.Sub(b.last).Seconds()
	if elapsed > 0 {
		b.tokens += elapsed * rate
		if b.tokens > capacity {
			b.tokens = capacity
		}
		b.last = now
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
