package domain

import "time"

// Scope grants access to a subset of tools.
type Scope string

const (
	ScopeRead         Scope = "read"
	ScopeWritePreview Scope = "write_preview"
	ScopeWriteExecute Scope = "write_execute"
	ScopeAdmin        Scope = "admin"
)

// AuthToken is an MCP authentication token (metadata only, never the secret).
type AuthToken struct {
	ID        string
	Hash      string
	Scopes    []Scope
	ExpiresAt time.Time
}

// HasScope reports whether the token grants a scope.
// ScopeAdmin implies every scope.
func (t AuthToken) HasScope(s Scope) bool {
	for _, sc := range t.Scopes {
		if sc == s || sc == ScopeAdmin {
			return true
		}
	}
	return false
}

// Expired reports whether the token is past its expiry.
// A zero ExpiresAt means the token never expires.
func (t AuthToken) Expired(now time.Time) bool {
	return !t.ExpiresAt.IsZero() && now.After(t.ExpiresAt)
}
