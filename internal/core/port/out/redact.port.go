package out

// Redactor removes secrets and PII from logs, audit records and responses.
// Implemented by adapter/out/redact/*.
type Redactor interface {
	RedactString(s string) string
	RedactValue(column, value string) string
}
