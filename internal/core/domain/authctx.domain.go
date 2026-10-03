package domain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type ctxKey string

const (
	scopesKey  ctxKey = "scopes"
	tokenIDKey ctxKey = "token_id"
	requestKey ctxKey = "request_id"
)

// WithAuth attaches the verified token identity to the context.
// adapter/in calls this after verifying the token; service reads it back.
func WithAuth(ctx context.Context, tokenID string, scopes []Scope) context.Context {
	ctx = context.WithValue(ctx, tokenIDKey, tokenID)
	return context.WithValue(ctx, scopesKey, scopes)
}

// ScopesFrom returns the scopes attached to the context.
func ScopesFrom(ctx context.Context) []Scope {
	sc, _ := ctx.Value(scopesKey).([]Scope)
	return sc
}

// TokenIDFrom returns the token id attached to the context.
func TokenIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(tokenIDKey).(string)
	return id
}

// WithRequestID attaches a request id to the context.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestKey, id)
}

// RequestIDFrom returns the request id, generating one when missing.
func RequestIDFrom(ctx context.Context) string {
	if id, _ := ctx.Value(requestKey).(string); id != "" {
		return id
	}
	return NewRequestID()
}

// NewRequestID generates a random request id.
func NewRequestID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return "req_" + hex.EncodeToString(b[:])
}

// RequireScope enforces token scopes on a context; ScopeAdmin implies all.
// Single source of truth for scope checks (service Gateway delegates here).
func RequireScope(ctx context.Context, s Scope) error {
	scopes := ScopesFrom(ctx)
	if len(scopes) == 0 {
		return NewError(CodeUnauthenticated, "missing authentication")
	}
	for _, sc := range scopes {
		if sc == s || sc == ScopeAdmin {
			return nil
		}
	}
	return NewError(CodeForbidden, "missing required scope")
}
