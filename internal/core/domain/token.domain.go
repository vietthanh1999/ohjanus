package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// TokenState tracks the lifecycle of a single-use write token.
type TokenState string

const (
	TokenPending  TokenState = "pending"
	TokenApproved TokenState = "approved"
	TokenRejected TokenState = "rejected"
	TokenUsed     TokenState = "used"
	TokenExpired  TokenState = "expired"
)

// PreviewToken binds a connection + SQL + params hash to a single execution.
// Any change to the SQL or params invalidates the token (TOKEN_MISMATCH).
// SQL/Params/Plan are kept so the Admin UI can display the pending request;
// params may contain PII and are never written to the audit log.
type PreviewToken struct {
	ID         string
	Connection string
	SQLHash    string
	ParamsHash string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	State      TokenState
	// ApprovedBy doubles as the decider id for rejections.
	ApprovedBy    string
	DecidedReason string
	// Display fields for the approval UI.
	SQL         string
	Params      []any
	PlanText    string
	Warnings    []string
	Statement   StatementType
	RequestedBy string
	Estimate    int64
	DecidedAt   time.Time
}

// TokenEventType describes token lifecycle transitions.
type TokenEventType string

const (
	TokenEventCreated  TokenEventType = "created"
	TokenEventApproved TokenEventType = "approved"
	TokenEventRejected TokenEventType = "rejected"
	TokenEventUsed     TokenEventType = "used"
	TokenEventExpired  TokenEventType = "expired"
)

// TokenEvent notifies subscribers (Admin SSE) of token transitions.
type TokenEvent struct {
	Type  TokenEventType
	Token *PreviewToken
}

// HashSQL hashes normalized SQL for token binding and audit records.
func HashSQL(sql string) string {
	sum := sha256.Sum256([]byte(sql))
	return "sha256:" + hex.EncodeToString(sum[:])
}

// HashParams hashes query params without revealing their values.
// Nil and empty params hash identically.
func HashParams(params []any) string {
	if len(params) == 0 {
		params = nil
	}
	raw, err := json.Marshal(params)
	if err != nil {
		raw = []byte(fmt.Sprintf("%v", params))
	}
	sum := sha256.Sum256(raw)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// ExecuteRequest is the inbound request for db_write_execute.
type ExecuteRequest struct {
	Connection     string
	SQL            string
	Params         []any
	PreviewTokenID string
}
