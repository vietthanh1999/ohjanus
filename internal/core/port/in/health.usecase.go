package in

import "context"

// HealthReport describes process liveness and pool readiness.
type HealthReport struct {
	Alive       bool
	Ready       bool
	Connections map[string]string
}

// HealthUseCase is the inbound port for health checks.
// Implemented by core/service, called by adapter/in.
type HealthUseCase interface {
	Health(ctx context.Context) HealthReport
}
