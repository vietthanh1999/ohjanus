package in

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// ReadUseCase is the inbound port for read-only queries.
// Implemented by core/service, called by adapter/in.
type ReadUseCase interface {
	Read(ctx context.Context, req domain.ReadRequest) (*domain.ResultSet, error)
	Explain(ctx context.Context, req domain.ReadRequest) (*domain.Plan, error)
}
