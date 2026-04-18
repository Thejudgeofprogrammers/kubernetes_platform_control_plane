package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"sync"
)

type MetricsInMemory struct {
	mu                     sync.RWMutex
	data                   map[string][]domain.Metric
	maxMetricsPerClientEnv int
}

func NewMetricsInMemory(
	maxMetricsPerClientEnv int,
) repository.MetricsRepository {
	return &MetricsInMemory{
		data:                   make(map[string][]domain.Metric),
		maxMetricsPerClientEnv: maxMetricsPerClientEnv,
	}
}

func (r *MetricsInMemory) Save(ctx context.Context, m domain.Metric) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	list := r.data[m.ClientID]
	list = append(list, m)

	if len(list) > r.maxMetricsPerClientEnv {
		list = list[len(list)-r.maxMetricsPerClientEnv:]
	}

	r.data[m.ClientID] = list

	return nil
}

func (r *MetricsInMemory) GetByClient(
	ctx context.Context,
	clientID string,
	limit int,
) ([]domain.Metric, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	metrics := r.data[clientID]

	if len(metrics) == 0 {
		return []domain.Metric{}, nil
	}

	if len(metrics) > limit {
		metrics = metrics[len(metrics)-limit:]
	}

	return metrics, nil
}

func (r *MetricsInMemory) DeleteByClient(
	ctx context.Context,
	clientID string,
) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, clientID)

	return nil
}
