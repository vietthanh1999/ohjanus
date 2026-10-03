package multi

import (
	"context"
	"errors"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Sink fans out audit events to every configured sink.
type Sink struct {
	sinks []out.AuditSink
}

var _ out.AuditSink = (*Sink)(nil)

// New builds a fan-out sink.
func New(sinks ...out.AuditSink) *Sink { return &Sink{sinks: sinks} }

// Emit forwards the event to all sinks.
func (s *Sink) Emit(ctx context.Context, e domain.AuditEvent) {
	for _, sn := range s.sinks {
		sn.Emit(ctx, e)
	}
}

// Close closes all sinks, joining errors.
func (s *Sink) Close() error {
	var errs []error
	for _, sn := range s.sinks {
		if err := sn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
