package out

import (
	"context"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// AuditFilter selects audit events for the Admin UI.
type AuditFilter struct {
	Event      string
	Connection string
	TokenID    string
	RequestID  string
	Status     string
	Search     string
	From       time.Time
	To         time.Time
	Limit      int
	Offset     int
}

// AuditPage is one page of audit events.
type AuditPage struct {
	Events []*domain.AuditEvent
	Total  int
	// NextOffset is -1 when there are no more pages.
	NextOffset int
}

// AuditReader queries past audit events for the Admin UI.
// Implemented by adapter/out/audit/* (queryable sinks only).
type AuditReader interface {
	Query(ctx context.Context, f AuditFilter) (AuditPage, error)
	Get(ctx context.Context, id string) (domain.AuditEvent, error)
}
