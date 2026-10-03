package config

import (
	"context"
	"crypto/subtle"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Resolver looks up MCP tokens from static config by SHA-256 hash.
type Resolver struct {
	tokens map[string]domain.AuthToken
	clock  out.Clock
}

var _ out.TokenResolver = (*Resolver)(nil)

// New builds a resolver from configured tokens.
func New(tokens []domain.AuthToken, clock out.Clock) *Resolver {
	m := make(map[string]domain.AuthToken, len(tokens))
	for _, t := range tokens {
		m[t.Hash] = t
	}
	return &Resolver{tokens: m, clock: clock}
}

// LookupByHash finds a token by hash using constant-time comparison.
func (r *Resolver) LookupByHash(_ context.Context, hash string) (domain.AuthToken, error) {
	for h, t := range r.tokens {
		if subtle.ConstantTimeCompare([]byte(h), []byte(hash)) == 1 {
			if t.Expired(r.clock.Now()) {
				return domain.AuthToken{}, domain.NewError(domain.CodeTokenExpired, "token expired")
			}
			return t, nil
		}
	}
	return domain.AuthToken{}, domain.NewError(domain.CodeTokenInvalid, "unknown token")
}
