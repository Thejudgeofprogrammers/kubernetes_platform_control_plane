package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"sync"
)

type InMemoryClientConfigRepository struct {
	storage map[string]*domain.APIClientConfig
	mu      sync.RWMutex
}

func NewInMemoryClientConfigRepository() repository.ClientConfigRepository {
	return &InMemoryClientConfigRepository{
		storage: make(map[string]*domain.APIClientConfig),
	}
}

func (r *InMemoryClientConfigRepository) Create(ctx context.Context, config *domain.APIClientConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[config.ID] = config
	return nil
}

func (r *InMemoryClientConfigRepository) GetByID(ctx context.Context, configID string) (*domain.APIClientConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	config, ok := r.storage[configID]
	if !ok {
		return nil, domain.ErrConfigNotFound
	}

	return config, nil
}

func (r *InMemoryClientConfigRepository) ListByClientID(ctx context.Context, clientID string) ([]*domain.APIClientConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.APIClientConfig, 0, len(r.storage))

	for _, config := range r.storage {
		if config.ClientID == clientID {
			result = append(result, config)
		}
	}

	return result, nil
}
