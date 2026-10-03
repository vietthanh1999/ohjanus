package sqlite

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite"

	"github.com/vietthanh1999/ohjanus/internal/adapter/out/auth/memory"
	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

const schema = `
CREATE TABLE IF NOT EXISTS auth_tokens (
	id         TEXT PRIMARY KEY,
	name       TEXT NOT NULL DEFAULT '',
	hash       TEXT NOT NULL UNIQUE,
	scopes     TEXT NOT NULL DEFAULT '[]',
	expires_at TEXT NOT NULL DEFAULT ''
);
`

// Store is a SQLite-backed AuthTokenStore. Runtime tokens survive restarts;
// config-file tokens are merged in via Seed (INSERT OR IGNORE, keyed by id
// and hash). Only hashes are stored — secrets are shown once at creation.
type Store struct {
	mu    sync.Mutex
	db    *sql.DB
	clock out.Clock
}

var _ out.AuthTokenStore = (*Store)(nil)

// Open creates the parent directory, enforces 0600 on the db file (it
// holds token hashes), and applies the schema.
func Open(path string, clock out.Clock) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("sqlite token store path is required")
	}
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, fmt.Errorf("create token store dir: %w", err)
		}
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("create token store file: %w", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("create token store file: %w", err)
	}
	if err := os.Chmod(path, 0o600); err != nil {
		return nil, fmt.Errorf("secure token store file: %w", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open token store: %w", err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA journal_mode=WAL;`); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable sqlite WAL: %w", err)
	}
	if _, err := db.Exec(`PRAGMA busy_timeout=5000;`); err != nil {
		db.Close()
		return nil, fmt.Errorf("set sqlite busy timeout: %w", err)
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate token store schema: %w", err)
	}
	return &Store{db: db, clock: clock}, nil
}

// Close releases the database handle.
func (s *Store) Close() error {
	return s.db.Close()
}

// Seed merges config-file tokens; conflicting ids or hashes keep the
// existing row so runtime state is never clobbered on restart.
func (s *Store) Seed(tokens []domain.AuthToken) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range tokens {
		scopes, err := json.Marshal(t.Scopes)
		if err != nil {
			return fmt.Errorf("encode token scopes: %w", err)
		}
		if _, err := s.db.Exec(
			`INSERT OR IGNORE INTO auth_tokens (id, name, hash, scopes, expires_at)
			 VALUES (?, ?, ?, ?, ?)`,
			t.ID, t.Name, t.Hash, string(scopes), formatTime(t.ExpiresAt),
		); err != nil {
			return fmt.Errorf("seed token %q: %w", t.ID, err)
		}
	}
	return nil
}

// LookupByHash finds a token by hash; expired tokens fail closed.
func (s *Store) LookupByHash(_ context.Context, hash string) (domain.AuthToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row := s.db.QueryRow(`SELECT id, name, hash, scopes, expires_at FROM auth_tokens WHERE hash = ?`, hash)
	t, err := scanToken(row)
	if err == sql.ErrNoRows {
		return domain.AuthToken{}, domain.NewError(domain.CodeTokenInvalid, "unknown token")
	}
	if err != nil {
		return domain.AuthToken{}, fmt.Errorf("lookup token: %w", err)
	}
	if t.Expired(s.clock.Now()) {
		return domain.AuthToken{}, domain.NewError(domain.CodeTokenExpired, "token expired")
	}
	return t, nil
}

// List returns token metadata sorted by id. Secrets are never stored.
func (s *Store) List() []domain.AuthToken {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.db.Query(`SELECT id, name, hash, scopes, expires_at FROM auth_tokens ORDER BY id`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []domain.AuthToken
	for rows.Next() {
		t, err := scanToken(rows)
		if err != nil {
			return nil
		}
		out = append(out, t)
	}
	return out
}

// Create generates a secret, stores only its hash, and returns the secret
// exactly once. Unknown scopes are rejected.
func (s *Store) Create(name string, scopes []domain.Scope, expiresAt time.Time) (string, domain.AuthToken, error) {
	if err := checkScopes(scopes); err != nil {
		return "", domain.AuthToken{}, err
	}
	secret, err := memory.GenerateSecret()
	if err != nil {
		return "", domain.AuthToken{}, err
	}
	t := domain.AuthToken{
		ID: newTokenID(), Name: name,
		Hash:      memory.HashSecret(secret),
		Scopes:    scopes,
		ExpiresAt: expiresAt,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	scopesJSON, err := json.Marshal(t.Scopes)
	if err != nil {
		return "", domain.AuthToken{}, fmt.Errorf("encode token scopes: %w", err)
	}
	if _, err := s.db.Exec(
		`INSERT INTO auth_tokens (id, name, hash, scopes, expires_at)
		 VALUES (?, ?, ?, ?, ?)`,
		t.ID, t.Name, t.Hash, string(scopesJSON), formatTime(t.ExpiresAt),
	); err != nil {
		return "", domain.AuthToken{}, fmt.Errorf("store token %q: %w", t.ID, err)
	}
	return secret, t, nil
}

// Revoke removes a token; subsequent uses fail with TOKEN_INVALID.
func (s *Store) Revoke(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	res, err := s.db.Exec(`DELETE FROM auth_tokens WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("revoke token %q: %w", id, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("revoke token %q: %w", id, err)
	}
	if n == 0 {
		return domain.NewError(domain.CodeTokenInvalid, fmt.Sprintf("unknown token %q", id))
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanToken(r rowScanner) (domain.AuthToken, error) {
	var t domain.AuthToken
	var scopesJSON, expiresAt string
	if err := r.Scan(&t.ID, &t.Name, &t.Hash, &scopesJSON, &expiresAt); err != nil {
		return domain.AuthToken{}, err
	}
	if err := json.Unmarshal([]byte(scopesJSON), &t.Scopes); err != nil {
		return domain.AuthToken{}, fmt.Errorf("decode token scopes: %w", err)
	}
	if expiresAt != "" {
		t.ExpiresAt, _ = time.Parse(time.RFC3339, expiresAt)
	}
	return t, nil
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func newTokenID() string {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("tok_%d", time.Now().UnixNano())
	}
	return "tok_" + hex.EncodeToString(b[:])
}

func checkScopes(scopes []domain.Scope) error {
	if len(scopes) == 0 {
		return fmt.Errorf("at least one scope is required")
	}
	for _, s := range scopes {
		valid := false
		for _, v := range memory.ValidScopes {
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
