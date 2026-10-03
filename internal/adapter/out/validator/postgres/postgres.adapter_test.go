package postgres

import (
	"context"
	"strings"
	"testing"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

func testValidator() *Validator {
	return New(map[string]ConnRules{
		"analytics": {
			AllowedSchemas: []string{"public"},
			DeniedTables:   []string{"users", "secrets"},
		},
		"open": {},
	}, nil)
}

func TestValidateReads(t *testing.T) {
	v := testValidator()
	ctx := context.Background()
	cases := []struct {
		name       string
		conn       string
		sql        string
		wantType   domain.StatementType
		wantTables []string
	}{
		{"simple", "open", "SELECT id, total FROM orders", domain.StatementSelect, []string{"orders"}},
		{"schema-qualified", "analytics", "SELECT * FROM public.orders", domain.StatementSelect, []string{"public.orders"}},
		{"alias-join-subquery", "open",
			"SELECT o.id FROM orders o JOIN (SELECT id FROM items) i ON i.id = o.id WHERE o.total > 10",
			domain.StatementSelect, []string{"orders", "items"}},
		{"cte", "open", "WITH recent AS (SELECT * FROM orders) SELECT * FROM recent", domain.StatementSelect, []string{"orders", "recent"}},
		{"comment-injection", "open", "SELECT/**/id/**/FROM/**/orders", domain.StatementSelect, []string{"orders"}},
		{"unicode", "open", "SELECT 'héllo' AS greeting", domain.StatementSelect, nil},
		{"no-from", "open", "SELECT 1", domain.StatementSelect, nil},
		{"explain", "open", "EXPLAIN SELECT * FROM orders", domain.StatementExplain, []string{"orders"}},
		{"show", "open", "SHOW transaction_isolation", domain.StatementShow, nil},
		{"update", "open", "UPDATE orders SET status = 'x'", domain.StatementUpdate, []string{"orders"}},
		{"insert", "open", "INSERT INTO orders (id) VALUES (1)", domain.StatementInsert, []string{"orders"}},
		{"delete", "open", "DELETE FROM orders WHERE id = 1", domain.StatementDelete, []string{"orders"}},
		{"modifying-cte", "open",
			"WITH m AS (UPDATE orders SET status = 'x' RETURNING *) SELECT * FROM m",
			domain.StatementUpdate, []string{"orders", "m"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vq, err := v.Validate(ctx, tc.conn, tc.sql)
			if err != nil {
				t.Fatalf("Validate: %v", err)
			}
			if vq.StatementType != tc.wantType {
				t.Errorf("type = %s, want %s", vq.StatementType, tc.wantType)
			}
			for _, want := range tc.wantTables {
				found := false
				for _, got := range vq.Tables {
					if strings.EqualFold(got, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("tables = %v, want %q included", vq.Tables, want)
				}
			}
			if vq.NormalizedSQL == "" {
				t.Error("normalized SQL is empty")
			}
		})
	}
}

func TestValidateDenials(t *testing.T) {
	v := testValidator()
	ctx := context.Background()
	cases := []struct {
		name     string
		conn     string
		sql      string
		wantCode domain.Code
	}{
		{"drop", "open", "DROP TABLE orders", domain.CodeQueryDenied},
		{"create", "open", "CREATE TABLE x (id int)", domain.CodeQueryDenied},
		{"truncate", "open", "TRUNCATE orders", domain.CodeQueryDenied},
		{"grant", "open", "GRANT SELECT ON orders TO reader", domain.CodeQueryDenied},
		{"pg-sleep", "open", "SELECT pg_sleep(1)", domain.CodeQueryDenied},
		{"pg-read-file", "open", "SELECT pg_read_file('x')", domain.CodeQueryDenied},
		{"schema-qualified-fn", "open", "SELECT pg_catalog.pg_sleep(1)", domain.CodeQueryDenied},
		{"denied-table", "analytics", "SELECT * FROM users", domain.CodeTableNotAllowed},
		{"denied-table-qualified", "analytics", "SELECT * FROM public.secrets", domain.CodeTableNotAllowed},
		{"denied-via-join", "analytics", "SELECT * FROM orders JOIN users ON users.id = orders.uid", domain.CodeTableNotAllowed},
		{"denied-via-subquery", "analytics", "SELECT * FROM (SELECT * FROM secrets) s", domain.CodeTableNotAllowed},
		{"schema-not-allowed", "analytics", "SELECT * FROM private.orders", domain.CodeSchemaNotAllowed},
		{"multi-statement", "open", "SELECT 1; SELECT 2", domain.CodeQueryDenied},
		{"empty", "open", "", domain.CodeQueryDenied},
		{"syntax-error", "open", "SELCT 1 FROM", domain.CodeParseError},
		{"explain-analyze", "open", "EXPLAIN ANALYZE SELECT * FROM orders", domain.CodeQueryDenied},
		{"copy-program", "open", "COPY (SELECT 1) TO PROGRAM 'id'", domain.CodeQueryDenied},
		{"do-block", "open", "DO $$ BEGIN RAISE NOTICE 'x'; END $$", domain.CodeQueryDenied},
		{"begin", "open", "BEGIN", domain.CodeQueryDenied},
		{"unknown-connection", "nope", "SELECT 1", domain.CodeConnectionNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := v.Validate(ctx, tc.conn, tc.sql)
			if err == nil {
				t.Fatalf("expected %s, got nil", tc.wantCode)
			}
			de, ok := err.(*domain.Error)
			if !ok {
				t.Fatalf("expected *domain.Error, got %T (%v)", err, err)
			}
			if de.Code != tc.wantCode {
				t.Errorf("code = %s, want %s (%s)", de.Code, tc.wantCode, de.Message)
			}
		})
	}
}

func TestValidateNormalized(t *testing.T) {
	v := testValidator()
	vq, err := v.Validate(context.Background(), "open", "select   a ,b   from   orders   where a>1")
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if !strings.Contains(vq.NormalizedSQL, "SELECT") {
		t.Errorf("normalized = %q, want SELECT in output", vq.NormalizedSQL)
	}
}
