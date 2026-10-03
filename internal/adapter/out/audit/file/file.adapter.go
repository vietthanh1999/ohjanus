package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Sink appends audit events as JSON lines to a file.
// Size-based rotation lands in v0.1 hardening (see rotate config).
type Sink struct {
	mu  sync.Mutex
	f   *os.File
	enc *json.Encoder
}

var _ out.AuditSink = (*Sink)(nil)

// New opens (or creates with 0600) the audit file for appending.
func New(path string) (*Sink, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open audit file: %w", err)
	}
	return &Sink{f: f, enc: json.NewEncoder(f)}, nil
}

// Emit appends one event.
func (s *Sink) Emit(_ context.Context, e domain.AuditEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.enc.Encode(e)
}

// Close closes the file.
func (s *Sink) Close() error { return s.f.Close() }
