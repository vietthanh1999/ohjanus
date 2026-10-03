package out

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// TokenStore persists single-use write tokens.
// In-memory for v0.1; Redis for multi-instance later.
// Implemented by adapter/out/token/*.
type TokenStore interface {
	Create(ctx context.Context, t *domain.PreviewToken) error
	Get(ctx context.Context, id string) (*domain.PreviewToken, error)
	Approve(ctx context.Context, id, decidedBy, reason string) (*domain.PreviewToken, error)
	Reject(ctx context.Context, id, decidedBy, reason string) (*domain.PreviewToken, error)
	MarkUsed(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.PreviewToken, error)
	// Watch streams lifecycle events (Admin SSE). Call cancel to unsubscribe.
	Watch() (<-chan domain.TokenEvent, func())
}
