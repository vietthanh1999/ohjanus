package memory

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

const tokenAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// ValidScopes are the scopes Janus understands.
var ValidScopes = []domain.Scope{
	domain.ScopeRead, domain.ScopeWritePreview, domain.ScopeWriteExecute, domain.ScopeAdmin,
}

// Store is an in-memory AuthTokenStore seeded from config.
// Tokens created at runtime are lost on restart; add them to the config
// file (via `janus token create`) to persist.
type Store struct {
	mu     sync.Mutex
	byID   map[string]*domain.AuthToken
	byHash map[string]string
	clock  out.Clock
}

var _ out.AuthTokenStore = (*Store)(nil)

// New builds an empty store.
func New(clock out.Clock) *Store {
	return &Store{byID: map[string]*domain.AuthToken{}, byHash: map[string]string{}, clock: clock}
}

// Seed loads config-file tokens into the store.
func (s *Store) Seed(tokens []domain.AuthToken) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range tokens {
		cp := t
		s.byID[cp.ID] = &cp
		s.byHash[cp.Hash] = cp.ID
	}
}

// LookupByHash finds a token by hash with constant-time comparison.
func (s *Store) LookupByHash(_ context.Context, hash string) (domain.AuthToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for h, id := range s.byHash {
		if subtle.ConstantTimeCompare([]byte(h), []byte(hash)) == 1 {
			t := s.byID[id]
			if t.Expired(s.clock.Now()) {
				return domain.AuthToken{}, domain.NewError(domain.CodeTokenExpired, "token expired")
			}
			return *t, nil
		}
	}
	return domain.AuthToken{}, domain.NewError(domain.CodeTokenInvalid, "unknown token")
}

// List returns token metadata sorted by id. Secrets are never stored.
func (s *Store) List() []domain.AuthToken {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]domain.AuthToken, 0, len(s.byID))
	for _, t := range s.byID {
		out = append(out, *t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Create generates a secret, stores only its hash, and returns the secret
// exactly once. Unknown scopes are rejected.
func (s *Store) Create(name string, scopes []domain.Scope, expiresAt time.Time) (string, domain.AuthToken, error) {
	if err := checkScopes(scopes); err != nil {
		return "", domain.AuthToken{}, err
	}
	secret, err := GenerateSecret()
	if err != nil {
		return "", domain.AuthToken{}, err
	}
	id, err := randomID("tok_")
	if err != nil {
		return "", domain.AuthToken{}, err
	}
	t := domain.AuthToken{ID: id, Name: name, Hash: hashSecret(secret), Scopes: scopes, ExpiresAt: expiresAt}
	s.mu.Lock()
	s.byID[t.ID] = &t
	s.byHash[t.Hash] = t.ID
	s.mu.Unlock()
	return secret, t, nil
}

// Revoke removes a token; subsequent uses fail with TOKEN_INVALID.
func (s *Store) Revoke(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.byID[id]
	if !ok {
		return domain.NewError(domain.CodeTokenInvalid, fmt.Sprintf("unknown token %q", id))
	}
	delete(s.byID, id)
	delete(s.byHash, t.Hash)
	return nil
}

func checkScopes(scopes []domain.Scope) error {
	if len(scopes) == 0 {
		return fmt.Errorf("at least one scope is required")
	}
	for _, s := range scopes {
		valid := false
		for _, v := range ValidScopes {
			if s == v {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("unknown scope %q", s)
		}
	}
	return nil
}

// GenerateSecret creates a random jn_ secret.
func GenerateSecret() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	sb := make([]byte, 32)
	for i := range sb {
		sb[i] = tokenAlphabet[int(b[i])%len(tokenAlphabet)]
	}
	return "jn_" + string(sb), nil
}

// HashSecret hashes a secret for storage and lookup.
func HashSecret(secret string) string {
	return hashSecret(secret)
}

func hashSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func randomID(prefix string) (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(b[:])[:8], nil
}
