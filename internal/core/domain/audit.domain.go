package domain

import "time"

// AuditEvent is a single structured audit record (JSONL).
// Params are never logged raw; only hashes are stored.
type AuditEvent struct {
	TS             time.Time     `json:"ts"`
	Event          string        `json:"event"`
	RequestID      string        `json:"request_id"`
	TokenID        string        `json:"token_id,omitempty"`
	Connection     string        `json:"connection,omitempty"`
	Tool           string        `json:"tool,omitempty"`
	SQLHash        string        `json:"sql_hash,omitempty"`
	SQLNormalized  string        `json:"sql_normalized,omitempty"`
	StatementType  StatementType `json:"statement_type,omitempty"`
	Tables         []string      `json:"tables,omitempty"`
	PolicyDecision string        `json:"policy_decision,omitempty"`
	PolicyRule     string        `json:"policy_rule,omitempty"`
	RowCount       int           `json:"row_count,omitempty"`
	Truncated      bool          `json:"truncated,omitempty"`
	DurationMs     int64         `json:"duration_ms,omitempty"`
	Status         string        `json:"status,omitempty"`
	Error          string        `json:"error,omitempty"`
}
