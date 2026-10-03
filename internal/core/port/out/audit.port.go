package out

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// AuditSink persists structured audit events.
// Implemented by adapter/out/audit/*.
type AuditSink interface {
	Emit(ctx context.Context, e domain.AuditEvent)
	Close() error
}
