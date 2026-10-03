package domain

import "time"

// TokenState tracks the lifecycle of a single-use write token.
type TokenState string

const (
	TokenPending  TokenState = "pending"
	TokenApproved TokenState = "approved"
	TokenUsed     TokenState = "used"
	TokenExpired  TokenState = "expired"
)

// PreviewToken binds a connection + SQL + params hash to a single execution.
// Any change to the SQL or params invalidates the token (TOKEN_MISMATCH).
type PreviewToken struct {
	ID         string
	Connection string
	SQLHash    string
	ParamsHash string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	State      TokenState
	ApprovedBy string
}

// ExecuteRequest is the inbound request for db_write_execute.
type ExecuteRequest struct {
	Connection     string
	SQL            string
	Params         []any
	PreviewTokenID string
}
