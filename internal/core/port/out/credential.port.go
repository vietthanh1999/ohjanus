package out

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// CredentialResolver resolves a DSN reference (env:/file:/keyring:/vault:)
// into credentials. Called once when a pool is opened.
// Implemented by adapter/out/credential/*.
type CredentialResolver interface {
	Scheme() string
	Resolve(ctx context.Context, ref string) (domain.Credentials, error)
}
