package denyall

import (
	"context"
	"fmt"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Validator denies every query. It is a fail-closed placeholder until the
// AST validator (pg_query_go) lands in Phase 1. It must never allow.
type Validator struct{}

var _ out.Validator = (*Validator)(nil)

// New builds the deny-all validator.
func New() *Validator { return &Validator{} }

// Validate always returns QUERY_DENIED.
func (v *Validator) Validate(_ context.Context, connection, _ string) (*domain.ValidatedQuery, error) {
	return nil, domain.ErrQueryDenied(
		fmt.Sprintf("no AST validator for connection %q yet (v0.1 scaffold denies all queries)", connection),
		"denyall",
	)
}
