package in

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// SchemaUseCase is the inbound port for schema introspection.
// Implemented by core/service, called by adapter/in.
type SchemaUseCase interface {
	ListConnections(ctx context.Context) ([]domain.Connection, error)
	GetSchema(ctx context.Context, connection, schema, table string) ([]domain.Schema, error)
}
