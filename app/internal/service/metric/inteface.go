package metric

import (
	"context"
	"control_plane/internal/domain"
)

type MetricsService interface {
	Collect(ctx context.Context, baseURL, clientID string) ([]domain.Metric, error)
}
