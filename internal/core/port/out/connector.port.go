package out

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// Connector opens a connection pool for one connection alias.
// Implemented by adapter/out/connector/*, one per driver.
type Connector interface {
	Driver() string
	Open(ctx context.Context, creds domain.Credentials, pool domain.PoolConfig) (Pool, error)
}

// Pool executes queries against one connection alias.
type Pool interface {
	Query(ctx context.Context, q domain.Query, opts domain.QueryOpts) (*domain.ResultSet, error)
	Exec(ctx context.Context, q domain.Query) (*domain.ExecResult, error)
	Explain(ctx context.Context, q domain.Query) (*domain.Plan, error)
	Schema(ctx context.Context, schema, table string) ([]domain.Schema, error)
	Ping(ctx context.Context) error
	Close() error
}
