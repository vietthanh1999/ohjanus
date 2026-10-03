package postgres

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// runPostgresContainer starts a throwaway Postgres and returns its DSN
// plus a cleanup func.
func runPostgresContainer(t *testing.T, ctx context.Context) (string, func()) {
	t.Helper()
	ctr, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("janus_test"),
		postgres.WithUsername("janus"),
		postgres.WithPassword("janus"),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = ctr.Terminate(ctx)
		t.Fatalf("connection string: %v", err)
	}
	return dsn, func() {
		if err := ctr.Terminate(context.Background()); err != nil {
			t.Logf("terminate container: %v", err)
		}
	}
}
