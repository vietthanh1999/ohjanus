package keyring

import (
	"context"
	"fmt"
	"strings"

	"github.com/zalando/go-keyring"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Resolver reads DSNs from the OS keychain (macOS Keychain, Linux Secret
// Service, Windows Credential Manager).
//
// Refs: keyring://service/key or keyring:key (service falls back to
// the configured default, "janus" when empty).
type Resolver struct {
	defaultService string
}

var _ out.CredentialResolver = (*Resolver)(nil)

// New builds a keyring resolver.
func New(defaultService string) *Resolver {
	if defaultService == "" {
		defaultService = "janus"
	}
	return &Resolver{defaultService: defaultService}
}

// Scheme returns "keyring".
func (r *Resolver) Scheme() string { return "keyring" }

// Resolve fetches the secret. The secret value is never logged.
func (r *Resolver) Resolve(_ context.Context, ref string) (domain.Credentials, error) {
	service, key, err := r.split(ref)
	if err != nil {
		return domain.Credentials{}, err
	}
	secret, err := keyring.Get(service, key)
	if err != nil {
		return domain.Credentials{}, fmt.Errorf("keyring %s/%s: %w", service, key, err)
	}
	if secret == "" {
		return domain.Credentials{}, fmt.Errorf("keyring %s/%s is empty", service, key)
	}
	return domain.Credentials{DSN: secret}, nil
}

func (r *Resolver) split(ref string) (service, key string, err error) {
	if rest, ok := strings.CutPrefix(ref, "keyring://"); ok {
		parts := strings.SplitN(rest, "/", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return "", "", fmt.Errorf("invalid keyring ref %q, want keyring://service/key", ref)
		}
		return parts[0], parts[1], nil
	}
	key, ok := strings.CutPrefix(ref, "keyring:")
	if !ok || key == "" {
		return "", "", fmt.Errorf("invalid keyring ref %q, want keyring://service/key or keyring:key", ref)
	}
	return r.defaultService, key, nil
}
