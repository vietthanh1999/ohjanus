package in

import (
	"context"
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
)

// WritePreview is the result of db_write_preview: an estimate plus a
// single-use token. Nothing is executed.
type WritePreview struct {
	PreviewTokenID   string
	ExpiresAt        time.Time
	StatementType    domain.StatementType
	AffectedEstimate int64
	Plan             string
	Warnings         []string
	SQLNormalized    string
}

// WriteResult is the result of db_write_execute.
type WriteResult struct {
	RowsAffected int64
	DurationMs   int64
}

// WriteUseCase is the inbound port for the write approval flow.
// Implemented by core/service, called by adapter/in.
type WriteUseCase interface {
	Preview(ctx context.Context, req domain.ReadRequest) (*WritePreview, error)
	Execute(ctx context.Context, req domain.ExecuteRequest) (*WriteResult, error)
}
