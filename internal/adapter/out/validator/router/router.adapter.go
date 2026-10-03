package router

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Router dispatches validation to the per-driver validator of each
// connection. Unknown connections fail closed via the fallback.
type Router struct {
	byConn   map[string]out.Validator
	fallback out.Validator
}

// New builds a router. A nil fallback denies unknown connections.
func New(byConn map[string]out.Validator, fallback out.Validator) *Router {
	if byConn == nil {
		byConn = map[string]out.Validator{}
	}
	return &Router{byConn: byConn, fallback: fallback}
}

// Validate routes to the connection's validator.
func (r *Router) Validate(ctx context.Context, connection, sql string) (*domain.ValidatedQuery, error) {
	if v, ok := r.byConn[connection]; ok {
		return v.Validate(ctx, connection, sql)
	}
	if r.fallback != nil {
		return r.fallback.Validate(ctx, connection, sql)
	}
	return nil, domain.ErrConnectionNotFound(connection)
}
