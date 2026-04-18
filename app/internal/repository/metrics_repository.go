package repository

import (
	"context"
	"control_plane/internal/domain"
)

type MetricsRepository interface {
	Save(ctx context.Context, m domain.Metric) error
	GetByClient(ctx context.Context, clientID string, limit int) ([]domain.Metric, error)
	DeleteByClient(ctx context.Context, clientID string) error
}
