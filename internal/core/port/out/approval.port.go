package out

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// ApprovalRequest asks a human to approve a write statement.
// PreviewTokenID correlates the request with its pending preview token so
// non-interactive engines (api) can wait for the Admin API/UI decision
// instead of prompting on a terminal.
type ApprovalRequest struct {
	Connection     string
	Statement      domain.StatementType
	SQL            string
	Params         []any
	Estimate       int64
	PreviewTokenID string
}

// ApprovalEngine obtains a human decision for a write.
// Implemented by adapter/out/approval/* (cli, http).
type ApprovalEngine interface {
	Approve(ctx context.Context, req ApprovalRequest) (bool, error)
}
