package service

import (
	"context"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/in"
)

// WriteService implements in.WriteUseCase.
// The full preview + approval + execute flow lands in v0.2;
// in v0.1 every call fails closed with APPROVAL_REQUIRED.
type WriteService struct {
	gateway *Gateway
}

var _ in.WriteUseCase = (*WriteService)(nil)

// NewWriteService wires a WriteService.
func NewWriteService(gateway *Gateway) *WriteService {
	return &WriteService{gateway: gateway}
}

// Preview refuses in v0.1: the approval engine does not exist yet.
func (s *WriteService) Preview(ctx context.Context, req domain.ReadRequest) (*in.WritePreview, error) {
	if err := s.gateway.RequireScope(ctx, domain.ScopeWritePreview); err != nil {
		return nil, err
	}
	return nil, domain.NewError(domain.CodeApprovalRequired, "write preview is not implemented in v0.1")
}

// Execute refuses in v0.1: the approval engine does not exist yet.
func (s *WriteService) Execute(ctx context.Context, req domain.ExecuteRequest) (*in.WriteResult, error) {
	if err := s.gateway.RequireScope(ctx, domain.ScopeWriteExecute); err != nil {
		return nil, err
	}
	return nil, domain.NewError(domain.CodeApprovalRequired, "write execute is not implemented in v0.1")
}
