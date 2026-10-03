package file

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Resolver reads DSNs from files: file:/path/to/secret.
// The file must have the expected mode (default 0600) or resolution fails.
type Resolver struct {
	expectedMode os.FileMode
}

var _ out.CredentialResolver = (*Resolver)(nil)

// New builds a file resolver. expectedMode 0 means 0600.
func New(expectedMode uint32) *Resolver {
	mode := os.FileMode(0o600)
	if expectedMode != 0 {
		mode = os.FileMode(expectedMode)
	}
	return &Resolver{expectedMode: mode}
}

// Scheme returns "file".
func (r *Resolver) Scheme() string { return "file" }

// Resolve reads and trims the secret file after checking its mode.
func (r *Resolver) Resolve(_ context.Context, ref string) (domain.Credentials, error) {
	path := strings.TrimPrefix(ref, "file:")
	if path == "" || path == ref {
		return domain.Credentials{}, fmt.Errorf("invalid file ref %q, want file:/path", ref)
	}
	info, err := os.Stat(path)
	if err != nil {
		return domain.Credentials{}, fmt.Errorf("stat secret file: %w", err)
	}
	if perm := info.Mode().Perm(); perm != r.expectedMode {
		return domain.Credentials{}, fmt.Errorf("secret file has mode %o, want %o", perm, r.expectedMode)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return domain.Credentials{}, fmt.Errorf("read secret file: %w", err)
	}
	if dsn := strings.TrimSpace(string(data)); dsn != "" {
		return domain.Credentials{DSN: dsn}, nil
	}
	return domain.Credentials{}, fmt.Errorf("secret file is empty")
}
