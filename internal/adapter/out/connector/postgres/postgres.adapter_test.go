package postgres

import (
	"context"
	"testing"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

func startPostgres(t *testing.T) (context.Context, string) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	ctx := context.Background()
	dsn, stop := runPostgresContainer(t, ctx)
	t.Cleanup(stop)
	return ctx, dsn
}

func TestPostgresQueryIntegration(t *testing.T) {
	ctx, dsn := startPostgres(t)
	pool, err := New().Open(ctx, domain.Credentials{DSN: dsn}, domain.PoolConfig{MaxOpen: 5})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer pool.Close()

	raw := pool.(*Pool).pool
	if _, err := raw.Exec(ctx, `CREATE TABLE orders (id bigint PRIMARY KEY, total numeric)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	for _, q := range []string{
		`INSERT INTO orders (id, total) VALUES (1, 100.5)`,
		`INSERT INTO orders (id, total) VALUES (2, 200)`,
		`INSERT INTO orders (id, total) VALUES (3, 300)`,
	} {
		if _, err := raw.Exec(ctx, q); err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	res, err := pool.Query(ctx, domain.Query{SQL: "SELECT id, total FROM orders ORDER BY id", Limit: 2},
		domain.QueryOpts{ReadOnly: true, RowLimit: 1000})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if res.RowCount != 2 || !res.Truncated {
		t.Errorf("got count=%d truncated=%v, want 2 true", res.RowCount, res.Truncated)
	}
	if len(res.Columns) != 2 || res.Columns[0] != "id" {
		t.Errorf("columns = %v", res.Columns)
	}

	res, err = pool.Query(ctx, domain.Query{SQL: "SELECT id FROM orders", Limit: 100},
		domain.QueryOpts{ReadOnly: true, RowLimit: 1000})
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if res.RowCount != 3 || res.Truncated {
		t.Errorf("got count=%d truncated=%v, want 3 false", res.RowCount, res.Truncated)
	}

	if _, err := pool.Query(ctx, domain.Query{SQL: "UPDATE orders SET total = 0"},
		domain.QueryOpts{ReadOnly: true}); err == nil {
		t.Error("expected read-only transaction to reject UPDATE, got nil")
	}
}

func TestPostgresExplainSchemaIntegration(t *testing.T) {
	ctx, dsn := startPostgres(t)
	pool, err := New().Open(ctx, domain.Credentials{DSN: dsn}, domain.PoolConfig{})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer pool.Close()

	raw := pool.(*Pool).pool
	if _, err := raw.Exec(ctx, `CREATE TABLE orders (id bigint PRIMARY KEY, total numeric)`); err != nil {
		t.Fatalf("create table: %v", err)
	}

	plan, err := pool.Explain(ctx, domain.Query{SQL: "SELECT * FROM orders"})
	if err != nil {
		t.Fatalf("Explain: %v", err)
	}
	if plan.Text == "" {
		t.Error("empty plan")
	}

	schemas, err := pool.Schema(ctx, "public", "orders")
	if err != nil {
		t.Fatalf("Schema: %v", err)
	}
	if len(schemas) != 1 || schemas[0].Name != "public" {
		t.Fatalf("schemas = %+v", schemas)
	}
	if len(schemas[0].Tables) != 1 || len(schemas[0].Tables[0].Columns) != 2 {
		t.Fatalf("tables = %+v", schemas[0].Tables)
	}
	if len(schemas[0].Tables[0].PrimaryKey) != 1 || schemas[0].Tables[0].PrimaryKey[0] != "id" {
		t.Errorf("pk = %v", schemas[0].Tables[0].PrimaryKey)
	}

	if err := pool.Ping(ctx); err != nil {
		t.Errorf("Ping: %v", err)
	}
}
