package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Connector opens pgx connection pools.
type Connector struct{}

var _ out.Connector = (*Connector)(nil)

// New builds a Postgres connector.
func New() *Connector { return &Connector{} }

// Driver returns "postgres".
func (c *Connector) Driver() string { return "postgres" }

// Open parses the DSN, applies pool settings and pings the database.
func (c *Connector) Open(ctx context.Context, creds domain.Credentials, pool domain.PoolConfig) (out.Pool, error) {
	if creds.DSN == "" {
		return nil, domain.NewError(domain.CodeDBError, "empty DSN")
	}
	cfg, err := pgxpool.ParseConfig(creds.DSN)
	if err != nil {
		return nil, domain.NewError(domain.CodeDBError, "invalid DSN")
	}
	if pool.MaxOpen > 0 {
		cfg.MaxConns = int32(pool.MaxOpen)
	}
	if pool.MaxIdle > 0 {
		cfg.MinConns = int32(pool.MaxIdle)
	}
	if pool.ConnMaxLifetime > 0 {
		cfg.MaxConnLifetime = pool.ConnMaxLifetime
	}
	p, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, dbErr(err)
	}
	// Retry the initial ping: the database may still be starting
	// (container boot, failover) when the gateway comes up.
	var pingErr error
	for i := 0; i < 15; i++ {
		pingErr = p.Ping(ctx)
		if pingErr == nil {
			return &Pool{pool: p}, nil
		}
		select {
		case <-ctx.Done():
			p.Close()
			return nil, dbErr(ctx.Err())
		case <-time.After(time.Second):
		}
	}
	p.Close()
	return nil, dbErr(pingErr)
}

// Pool executes queries against one Postgres database.
type Pool struct {
	pool *pgxpool.Pool
}

var _ out.Pool = (*Pool)(nil)

// Query runs a read query with a hard LIMIT wrap (§6.5).
// It fetches limit+1 rows to detect truncation.
func (p *Pool) Query(ctx context.Context, q domain.Query, opts domain.QueryOpts) (*domain.ResultSet, error) {
	limit := q.Limit
	if limit <= 0 {
		limit = opts.RowLimit
	}
	if limit <= 0 {
		limit = 1000
	}
	wrapped := fmt.Sprintf("SELECT * FROM (%s) AS _janus_limit LIMIT %d", trimTrailingSemi(q.SQL), limit+1)
	start := time.Now()
	query := p.pool.Query
	if opts.ReadOnly {
		var res *domain.ResultSet
		err := withReadOnlyTx(ctx, p.pool, func(tx pgx.Tx) error {
			r, rerr := collect(ctx, tx.Query, wrapped, q.Params, limit, start)
			if rerr != nil {
				return rerr
			}
			res = r
			return nil
		})
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return collect(ctx, query, wrapped, q.Params, limit, start)
}

type queryFunc func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

// Exec runs a write statement in a read-write transaction.
func (p *Pool) Exec(ctx context.Context, q domain.Query) (*domain.ExecResult, error) {
	start := time.Now()
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, dbErr(err)
	}
	defer tx.Rollback(ctx)
	tag, err := tx.Exec(ctx, q.SQL, q.Params...)
	if err != nil {
		return nil, dbErr(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, dbErr(err)
	}
	return &domain.ExecResult{RowsAffected: tag.RowsAffected(), DurationMs: time.Since(start).Milliseconds()}, nil
}

func collect(ctx context.Context, query queryFunc, sql string, params []any, limit int, start time.Time) (*domain.ResultSet, error) {
	rows, err := query(ctx, sql, params...)
	if err != nil {
		return nil, dbErr(err)
	}
	defer rows.Close()
	cols := make([]string, 0, len(rows.FieldDescriptions()))
	for _, fd := range rows.FieldDescriptions() {
		cols = append(cols, fd.Name)
	}
	var out [][]any
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, dbErr(err)
		}
		out = append(out, vals)
	}
	if err := rows.Err(); err != nil {
		return nil, dbErr(err)
	}
	truncated := false
	if len(out) > limit {
		out = out[:limit]
		truncated = true
	}
	return &domain.ResultSet{
		Columns: cols, Rows: out, RowCount: len(out),
		Truncated: truncated, DurationMs: time.Since(start).Milliseconds(),
	}, nil
}

// Explain runs EXPLAIN without executing the query.
func (p *Pool) Explain(ctx context.Context, q domain.Query) (*domain.Plan, error) {
	rows, err := p.pool.Query(ctx, "EXPLAIN "+q.SQL, q.Params...)
	if err != nil {
		return nil, dbErr(err)
	}
	defer rows.Close()
	var lines []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return nil, dbErr(err)
		}
		lines = append(lines, line)
	}
	if err := rows.Err(); err != nil {
		return nil, dbErr(err)
	}
	return &domain.Plan{Text: strings.Join(lines, "\n")}, nil
}

// Schema introspects tables and columns via information_schema.
func (p *Pool) Schema(ctx context.Context, schema, table string) ([]domain.Schema, error) {
	colsQuery := `SELECT table_schema, table_name, column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')`
	pkQuery := `SELECT tc.table_schema, tc.table_name, kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
		  ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema
		WHERE tc.constraint_type = 'PRIMARY KEY'
		  AND tc.table_schema NOT IN ('pg_catalog', 'information_schema')`
	var args []any
	var pkArgs []any
	if schema != "" {
		colsQuery += fmt.Sprintf(" AND table_schema = $%d", len(args)+1)
		pkQuery += fmt.Sprintf(" AND tc.table_schema = $%d", len(pkArgs)+1)
		args = append(args, schema)
		pkArgs = append(pkArgs, schema)
	}
	if table != "" {
		colsQuery += fmt.Sprintf(" AND table_name = $%d", len(args)+1)
		pkQuery += fmt.Sprintf(" AND tc.table_name = $%d", len(pkArgs)+1)
		args = append(args, table)
		pkArgs = append(pkArgs, table)
	}
	colsQuery += " ORDER BY table_schema, table_name, ordinal_position"

	type col struct {
		schema, table, name, typ string
		nullable                 bool
	}
	var columns []col
	rows, err := p.pool.Query(ctx, colsQuery, args...)
	if err != nil {
		return nil, dbErr(err)
	}
	for rows.Next() {
		var c col
		var nullable string
		if err := rows.Scan(&c.schema, &c.table, &c.name, &c.typ, &nullable); err != nil {
			rows.Close()
			return nil, dbErr(err)
		}
		c.nullable = nullable == "YES"
		columns = append(columns, c)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, dbErr(err)
	}
	rows.Close()

	pks := map[string][]string{}
	prows, err := p.pool.Query(ctx, pkQuery, pkArgs...)
	if err != nil {
		return nil, dbErr(err)
	}
	for prows.Next() {
		var sch, tbl, cn string
		if err := prows.Scan(&sch, &tbl, &cn); err != nil {
			prows.Close()
			return nil, dbErr(err)
		}
		pks[sch+"."+tbl] = append(pks[sch+"."+tbl], cn)
	}
	if err := prows.Err(); err != nil {
		prows.Close()
		return nil, dbErr(err)
	}
	prows.Close()

	var schemas []domain.Schema
	idx := map[string]int{}
	for _, c := range columns {
		i, ok := idx[c.schema]
		if !ok {
			schemas = append(schemas, domain.Schema{Name: c.schema})
			i = len(schemas) - 1
			idx[c.schema] = i
		}
		ti := -1
		for j := range schemas[i].Tables {
			if schemas[i].Tables[j].Name == c.table {
				ti = j
				break
			}
		}
		if ti < 0 {
			schemas[i].Tables = append(schemas[i].Tables, domain.Table{
				Name:       c.table,
				PrimaryKey: pks[c.schema+"."+c.table],
			})
			ti = len(schemas[i].Tables) - 1
		}
		t := &schemas[i].Tables[ti]
		t.Columns = append(t.Columns, domain.Column{Name: c.name, Type: c.typ, Nullable: c.nullable})
	}
	return schemas, nil
}

// Ping checks the pool.
func (p *Pool) Ping(ctx context.Context) error {
	if err := p.pool.Ping(ctx); err != nil {
		return dbErr(err)
	}
	return nil
}

// Close drains the pool.
func (p *Pool) Close() error {
	p.pool.Close()
	return nil
}

// withReadOnlyTx runs fn inside SET TRANSACTION READ ONLY (§6.4).
func withReadOnlyTx(ctx context.Context, pool *pgxpool.Pool, fn func(pgx.Tx) error) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{AccessMode: pgx.ReadOnly})
	if err != nil {
		return dbErr(err)
	}
	defer tx.Rollback(ctx)
	if err := fn(tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return dbErr(err)
	}
	return nil
}

func dbErr(err error) error {
	return domain.NewError(domain.CodeDBError, err.Error())
}

func trimTrailingSemi(sql string) string {
	return strings.TrimSuffix(strings.TrimSpace(sql), ";")
}
