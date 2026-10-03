package env

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Resolver reads DSNs from environment variables: env:VAR_NAME.
type Resolver struct{}

var _ out.CredentialResolver = (*Resolver)(nil)

// New builds an env resolver.
func New() *Resolver { return &Resolver{} }

// Scheme returns "env".
func (r *Resolver) Scheme() string { return "env" }

// Resolve looks up VAR_NAME in the process environment.
func (r *Resolver) Resolve(_ context.Context, ref string) (domain.Credentials, error) {
	name := strings.TrimPrefix(ref, "env:")
	if name == "" || name == ref {
		return domain.Credentials{}, fmt.Errorf("invalid env ref %q, want env:VAR_NAME", ref)
	}
	val, ok := os.LookupEnv(name)
	if !ok || val == "" {
		return domain.Credentials{}, fmt.Errorf("environment variable %q is not set", name)
	}
	return domain.Credentials{DSN: val}, nil
}
