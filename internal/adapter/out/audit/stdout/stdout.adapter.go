package stdout

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Sink writes audit events as JSON lines to stdout.
type Sink struct {
	mu  sync.Mutex
	enc *json.Encoder
}

var _ out.AuditSink = (*Sink)(nil)

// New builds a stdout sink.
func New() *Sink { return &Sink{enc: json.NewEncoder(os.Stdout)} }

// Emit writes one event.
func (s *Sink) Emit(_ context.Context, e domain.AuditEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.enc.Encode(e)
}

// Close is a no-op.
func (s *Sink) Close() error { return nil }
