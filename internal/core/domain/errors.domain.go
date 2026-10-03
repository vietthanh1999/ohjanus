package domain

import "fmt"

// Code is a machine-readable error code (§12.1 of the spec).
type Code string

const (
	CodeUnauthenticated    Code = "UNAUTHENTICATED"
	CodeTokenInvalid       Code = "TOKEN_INVALID"
	CodeTokenExpired       Code = "TOKEN_EXPIRED"
	CodeForbidden          Code = "FORBIDDEN"
	CodeConnectionNotFound Code = "CONNECTION_NOT_FOUND"
	CodeSchemaNotAllowed   Code = "SCHEMA_NOT_ALLOWED"
	CodeTableNotAllowed    Code = "TABLE_NOT_ALLOWED"
	CodeParseError         Code = "PARSE_ERROR"
	CodeQueryDenied        Code = "QUERY_DENIED"
	CodeQueryTooLong       Code = "QUERY_TOO_LONG"
	CodeQueryTimeout       Code = "QUERY_TIMEOUT"
	CodeRateLimited        Code = "RATE_LIMITED"
	CodeTooManyConcurrent  Code = "TOO_MANY_CONCURRENT_QUERIES"
	CodeTokenMismatch      Code = "TOKEN_MISMATCH"
	CodeTokenAlreadyUsed   Code = "TOKEN_ALREADY_USED"
	CodeApprovalRequired   Code = "APPROVAL_REQUIRED"
	CodeDBError            Code = "DB_ERROR"
	CodeInternal           Code = "INTERNAL"
)

// Error is a Janus error with a stable code for MCP responses.
type Error struct {
	Code      Code
	Message   string
	Rule      string
	RequestID string
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewError builds a plain coded error.
func NewError(code Code, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

// ErrQueryDenied builds a QUERY_DENIED error bound to a policy rule.
func ErrQueryDenied(reason, rule string) *Error {
	return &Error{Code: CodeQueryDenied, Message: reason, Rule: rule}
}

// ErrConnectionNotFound builds a CONNECTION_NOT_FOUND error.
func ErrConnectionNotFound(name string) *Error {
	return &Error{Code: CodeConnectionNotFound, Message: fmt.Sprintf("connection %q not found", name)}
}
