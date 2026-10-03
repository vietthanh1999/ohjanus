package chain

import (
	"context"
	"fmt"
	"strings"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Chain dispatches a DSN ref to the resolver registered for its scheme
// prefix (env:, file:, ...).
type Chain struct {
	resolvers map[string]out.CredentialResolver
}

// New builds a chain from the given resolvers.
func New(resolvers ...out.CredentialResolver) *Chain {
	m := make(map[string]out.CredentialResolver, len(resolvers))
	for _, r := range resolvers {
		m[r.Scheme()] = r
	}
	return &Chain{resolvers: m}
}

// Resolve routes ref to the matching scheme resolver.
func (c *Chain) Resolve(ctx context.Context, ref string) (domain.Credentials, error) {
	scheme := ref
	if i := strings.Index(ref, ":"); i >= 0 {
		scheme = ref[:i]
	}
	r, ok := c.resolvers[scheme]
	if !ok {
		return domain.Credentials{}, fmt.Errorf("no credential resolver for scheme %q", scheme)
	}
	return r.Resolve(ctx, ref)
}
