package domain

import "time"

// StatementType is the top-level SQL statement kind, derived from the AST.
type StatementType string

const (
	StatementSelect   StatementType = "SELECT"
	StatementWith     StatementType = "WITH"
	StatementExplain  StatementType = "EXPLAIN"
	StatementShow     StatementType = "SHOW"
	StatementInsert   StatementType = "INSERT"
	StatementUpdate   StatementType = "UPDATE"
	StatementDelete   StatementType = "DELETE"
	StatementMerge    StatementType = "MERGE"
	StatementCreate   StatementType = "CREATE"
	StatementAlter    StatementType = "ALTER"
	StatementDrop     StatementType = "DROP"
	StatementTruncate StatementType = "TRUNCATE"
	StatementGrant    StatementType = "GRANT"
	StatementRevoke   StatementType = "REVOKE"
	StatementUnknown  StatementType = "UNKNOWN"
)

// IsRead reports whether the statement type belongs to the read path.
func (s StatementType) IsRead() bool {
	switch s {
	case StatementSelect, StatementWith, StatementExplain, StatementShow:
		return true
	default:
		return false
	}
}

// IsWrite reports whether the statement type belongs to the write path
// (requires preview + approval).
func (s StatementType) IsWrite() bool {
	switch s {
	case StatementInsert, StatementUpdate, StatementDelete, StatementMerge:
		return true
	default:
		return false
	}
}

// Query is a parameterized query against a connection alias.
type Query struct {
	Connection string
	SQL        string
	Params     []any
	Limit      int
}

// ReadRequest is the inbound request for db_read / db_explain.
type ReadRequest struct {
	Connection string
	SQL        string
	Params     []any
	Limit      int
}

// ResultSet is the tabular result of a read query.
type ResultSet struct {
	Columns    []string
	Rows       [][]any
	RowCount   int
	Truncated  bool
	DurationMs int64
}

// ExecResult is the outcome of a write execution.
type ExecResult struct {
	RowsAffected int64
	DurationMs   int64
}

// ValidatedQuery is the output of AST validation.
type ValidatedQuery struct {
	Query         Query
	StatementType StatementType
	Tables        []string
	Schemas       []string
	Functions     []string
	NormalizedSQL string
}

// Plan is the output of EXPLAIN (never executes the query).
type Plan struct {
	Text             string
	AffectedEstimate int64
}

// QueryOpts controls how a query is executed by the pool.
type QueryOpts struct {
	ReadOnly bool
	Timeout  time.Duration
	RowLimit int
}
