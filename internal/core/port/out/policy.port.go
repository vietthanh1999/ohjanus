package out

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// PolicyEngine evaluates validated queries against ordered rules.
// Implemented by adapter/out/policy/*.
type PolicyEngine interface {
	Evaluate(ctx context.Context, q *domain.ValidatedQuery) domain.PolicyDecision
}
