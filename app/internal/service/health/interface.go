package health

import (
	"context"
	"control_plane/internal/domain"
)

type HealthService interface {
	Update(ctx context.Context, clientID string, status domain.HealthStatus, message string) error
	Get(ctx context.Context, clientID string) (*domain.APIClientHealth, error)
	Set(clientID string, status domain.HealthStatus)
}
