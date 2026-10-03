package out

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// Validator parses SQL into an AST and enforces table/schema/function rules.
// Implemented by adapter/out/validator/*, one per driver.
type Validator interface {
	Validate(ctx context.Context, connection, sql string) (*domain.ValidatedQuery, error)
}
