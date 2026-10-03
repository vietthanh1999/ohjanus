package vault

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	vault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Resolver reads DSNs from HashiCorp Vault KV v2 with AppRole auth.
//
// Refs: vault://mount/path/to/secret#field, e.g.
// vault://secret/data/janus/analytics#dsn
// (a "data/" segment after the mount is stripped for KV v2 URLs;
// field defaults to "dsn").
type Resolver struct {
	address          string
	roleID, secretID string

	mu     sync.Mutex
	client *vault.Client
}

var _ out.CredentialResolver = (*Resolver)(nil)

// New builds a Vault resolver. Login happens lazily on first Resolve.
func New(address, roleID, secretID string) (*Resolver, error) {
	if address == "" {
		return nil, fmt.Errorf("vault address is required")
	}
	if roleID == "" || secretID == "" {
		return nil, fmt.Errorf("vault approle role_id and secret_id are required")
	}
	return &Resolver{address: address, roleID: roleID, secretID: secretID}, nil
}

// Scheme returns "vault".
func (r *Resolver) Scheme() string { return "vault" }

// Resolve fetches the field. The secret value is never logged.
func (r *Resolver) Resolve(ctx context.Context, ref string) (domain.Credentials, error) {
	mount, path, field, err := splitRef(ref)
	if err != nil {
		return domain.Credentials{}, err
	}
	val, err := r.read(ctx, mount, path, field)
	if err != nil {
		if isForbidden(err) {
			// Token may have expired: re-login once and retry.
			if lerr := r.login(ctx); lerr != nil {
				return domain.Credentials{}, lerr
			}
			val, err = r.read(ctx, mount, path, field)
		}
		if err != nil {
			return domain.Credentials{}, err
		}
	}
	if val == "" {
		return domain.Credentials{}, fmt.Errorf("vault %s/%s has empty field %q", mount, path, field)
	}
	return domain.Credentials{DSN: val}, nil
}

func (r *Resolver) read(ctx context.Context, mount, path, field string) (string, error) {
	client, err := r.ensureClient(ctx)
	if err != nil {
		return "", err
	}
	resp, err := client.Secrets.KvV2Read(ctx, path, vault.WithMountPath(mount))
	if err != nil {
		return "", fmt.Errorf("vault read %s/%s: %w", mount, path, err)
	}
	raw, ok := resp.Data.Data[field]
	if !ok {
		return "", fmt.Errorf("vault %s/%s missing field %q", mount, path, field)
	}
	str, ok := raw.(string)
	if !ok || str == "" {
		return "", fmt.Errorf("vault %s/%s field %q is not a non-empty string", mount, path, field)
	}
	return str, nil
}

func (r *Resolver) ensureClient(ctx context.Context) (*vault.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.client != nil {
		return r.client, nil
	}
	return r.loginLocked(ctx)
}

func (r *Resolver) login(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.loginLocked(ctx)
	return err
}

func (r *Resolver) loginLocked(ctx context.Context) (*vault.Client, error) {
	client, err := vault.New(vault.WithAddress(r.address))
	if err != nil {
		return nil, fmt.Errorf("vault client: %w", err)
	}
	resp, err := client.Auth.AppRoleLogin(ctx, schema.AppRoleLoginRequest{
		RoleId: r.roleID, SecretId: r.secretID,
	})
	if err != nil {
		return nil, fmt.Errorf("vault approle login: %w", err)
	}
	if resp.Auth == nil || resp.Auth.ClientToken == "" {
		return nil, fmt.Errorf("vault approle login returned no token")
	}
	if err := client.SetToken(resp.Auth.ClientToken); err != nil {
		return nil, fmt.Errorf("vault set token: %w", err)
	}
	r.client = client
	return client, nil
}

func isForbidden(err error) bool {
	var re *vault.ResponseError
	return errors.As(err, &re) && re.StatusCode == http.StatusForbidden
}

func splitRef(ref string) (mount, path, field string, err error) {
	rest, ok := strings.CutPrefix(ref, "vault://")
	if !ok || rest == "" {
		return "", "", "", fmt.Errorf("invalid vault ref %q, want vault://mount/path#field", ref)
	}
	rest, field, _ = strings.Cut(rest, "#")
	if field == "" {
		field = "dsn"
	}
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", "", fmt.Errorf("invalid vault ref %q, want vault://mount/path#field", ref)
	}
	path = strings.TrimPrefix(parts[1], "data/")
	if path == "" {
		return "", "", "", fmt.Errorf("invalid vault ref %q, empty secret path", ref)
	}
	return parts[0], path, field, nil
}
