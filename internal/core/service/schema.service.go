package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/in"
)

// SchemaService implements in.SchemaUseCase with allow/deny filtering.
type SchemaService struct {
	gateway *Gateway
	conns   *ConnRegistry
}

var _ in.SchemaUseCase = (*SchemaService)(nil)

// NewSchemaService wires a SchemaService.
func NewSchemaService(gateway *Gateway, conns *ConnRegistry) *SchemaService {
	if conns == nil {
		conns = NewConnRegistry(nil, nil)
	}
	return &SchemaService{gateway: gateway, conns: conns}
}

// ListConnections returns the connection aliases the agent may see.
// Credentials are never included.
func (s *SchemaService) ListConnections(ctx context.Context) ([]domain.Connection, error) {
	status := "success"
	defer func() { s.gateway.Observe("db_list_connections", status) }()
	if err := s.gateway.RequireScope(ctx, domain.ScopeRead); err != nil {
		status = "error"
		return nil, err
	}
	out := make([]domain.Connection, 0)
	for _, m := range s.conns.Metas() {
		out = append(out, m.Connection)
	}
	return out, nil
}

// GetSchema returns filtered schemas/tables for a connection.
func (s *SchemaService) GetSchema(ctx context.Context, connection, schema, table string) ([]domain.Schema, error) {
	status := "success"
	defer func() { s.gateway.Observe("db_schema", status) }()
	if err := s.gateway.RequireScope(ctx, domain.ScopeRead); err != nil {
		status = "error"
		return nil, err
	}
	meta, ok := s.conns.Meta(connection)
	if !ok {
		status = "error"
		return nil, domain.ErrConnectionNotFound(connection)
	}
	if schema != "" && len(meta.AllowedSchemas) > 0 && !containsFold(meta.AllowedSchemas, schema) {
		status = "denied"
		return nil, domain.NewError(domain.CodeSchemaNotAllowed,
			fmt.Sprintf("schema %q is not allowed", schema))
	}
	if table != "" && isTableDenied(meta, schema, table) {
		status = "denied"
		return nil, domain.NewError(domain.CodeTableNotAllowed,
			fmt.Sprintf("table %q is not allowed", table))
	}
	pool, ok := s.conns.Pool(connection)
	if !ok {
		status = "error"
		return nil, domain.ErrConnectionNotFound(connection)
	}
	schemas, err := pool.Schema(ctx, schema, table)
	if err != nil {
		status = "error"
		return nil, err
	}
	return filterSchemas(schemas, meta), nil
}

// filterSchemas keeps only allowed schemas/tables.
// Empty AllowedSchemas / AllowedTables means "all except denied".
func filterSchemas(schemas []domain.Schema, meta domain.ConnectionMeta) []domain.Schema {
	out := make([]domain.Schema, 0, len(schemas))
	for _, sc := range schemas {
		if len(meta.AllowedSchemas) > 0 && !containsFold(meta.AllowedSchemas, sc.Name) {
			continue
		}
		kept := make([]domain.Table, 0, len(sc.Tables))
		for _, t := range sc.Tables {
			if isTableDenied(meta, sc.Name, t.Name) {
				continue
			}
			kept = append(kept, t)
		}
		sc.Tables = kept
		out = append(out, sc)
	}
	return out
}

func isTableDenied(meta domain.ConnectionMeta, schema, table string) bool {
	qualified := table
	if schema != "" {
		qualified = schema + "." + table
	}
	for _, d := range meta.DeniedTables {
		if strings.EqualFold(d, table) || strings.EqualFold(d, qualified) {
			return true
		}
	}
	if len(meta.AllowedTables) == 0 {
		return false
	}
	for _, a := range meta.AllowedTables {
		if strings.EqualFold(a, table) || strings.EqualFold(a, qualified) {
			return false
		}
	}
	return true
}

func containsFold(list []string, s string) bool {
	for _, v := range list {
		if strings.EqualFold(v, s) {
			return true
		}
	}
	return false
}
