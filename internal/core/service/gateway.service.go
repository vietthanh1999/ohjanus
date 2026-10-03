package service

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Gateway enforces cross-cutting concerns: scopes, limits, timeouts.
// Auth identity travels on the context via domain.WithAuth (adapter/in sets
// it after verifying the token), so adapter/in never imports this package.
type Gateway struct {
	clock out.Clock
}

// NewGateway builds a Gateway.
func NewGateway(clock out.Clock) *Gateway {
	return &Gateway{clock: clock}
}

// RequireScope enforces token scopes; ScopeAdmin implies every scope.
func (g *Gateway) RequireScope(ctx context.Context, s domain.Scope) error {
	return domain.RequireScope(ctx, s)
}
