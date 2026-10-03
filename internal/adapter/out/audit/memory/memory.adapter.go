package memory

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// defaultCapacity bounds memory use of the queryable audit trail.
// The file sink remains the durable record; this buffer serves the UI.
const defaultCapacity = 10000

// Buffer is a queryable in-memory audit sink (ring of the last N events).
// Wire it alongside durable sinks via audit/multi.
type Buffer struct {
	mu     sync.Mutex
	events []*domain.AuditEvent
	cap    int
	seq    atomic.Uint64
}

var (
	_ out.AuditSink   = (*Buffer)(nil)
	_ out.AuditReader = (*Buffer)(nil)
)

// New builds a buffer holding up to capacity events (<=0 means default).
func New(capacity int) *Buffer {
	if capacity <= 0 {
		capacity = defaultCapacity
	}
	return &Buffer{cap: capacity}
}

// Emit assigns an id and appends the event, dropping the oldest on overflow.
func (b *Buffer) Emit(_ context.Context, e domain.AuditEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if e.ID == "" {
		e.ID = fmt.Sprintf("aud_%d", b.seq.Add(1))
	}
	cp := e
	b.events = append(b.events, &cp)
	if len(b.events) > b.cap {
		kept := make([]*domain.AuditEvent, b.cap)
		copy(kept, b.events[len(b.events)-b.cap:])
		b.events = kept
	}
}

// Close is a no-op.
func (b *Buffer) Close() error { return nil }

// Get returns one event by id.
func (b *Buffer) Get(_ context.Context, id string) (domain.AuditEvent, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, e := range b.events {
		if e.ID == id {
			return *e, nil
		}
	}
	return domain.AuditEvent{}, fmt.Errorf("audit event %q not found", id)
}

// Query filters newest-first with offset pagination.
func (b *Buffer) Query(_ context.Context, f out.AuditFilter) (out.AuditPage, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	var matched []*domain.AuditEvent
	for i := len(b.events) - 1; i >= 0; i-- {
		if filterMatch(b.events[i], f) {
			matched = append(matched, b.events[i])
		}
	}
	total := len(matched)
	limit := f.Limit
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	page := out.AuditPage{Total: total, NextOffset: -1}
	for _, e := range matched[offset:end] {
		cp := *e
		page.Events = append(page.Events, &cp)
	}
	if end < total {
		page.NextOffset = end
	}
	return page, nil
}

func filterMatch(e *domain.AuditEvent, f out.AuditFilter) bool {
	if f.Event != "" && e.Event != f.Event {
		return false
	}
	if f.Connection != "" && e.Connection != f.Connection {
		return false
	}
	if f.TokenID != "" && e.TokenID != f.TokenID {
		return false
	}
	if f.RequestID != "" && e.RequestID != f.RequestID {
		return false
	}
	if f.Status != "" && e.Status != f.Status {
		return false
	}
	if !f.From.IsZero() && e.TS.Before(f.From) {
		return false
	}
	if !f.To.IsZero() && e.TS.After(f.To) {
		return false
	}
	if f.Search != "" && !strings.Contains(strings.ToLower(e.SQLNormalized), strings.ToLower(f.Search)) {
		return false
	}
	return true
}
