package out

import (
	"context"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// TokenResolver looks up an MCP auth token by its SHA-256 hash.
// Implemented by adapter/out/auth/*.
type TokenResolver interface {
	LookupByHash(ctx context.Context, hash string) (domain.AuthToken, error)
}

// AuthTokenStore manages MCP auth tokens at runtime.
// Tokens created here are runtime-only (lost on restart); persistent tokens
// live in the config file. Implemented by adapter/out/auth/*.
type AuthTokenStore interface {
	TokenResolver
	List() []domain.AuthToken
	Create(name string, scopes []domain.Scope, expiresAt time.Time) (secret string, token domain.AuthToken, err error)
	Revoke(id string) error
}
