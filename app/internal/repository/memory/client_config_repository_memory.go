package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/repository"
	"sync"
)

type InMemoryClientConfigRepository struct {
	storage map[string]*domain.APIClientConfig
	mu      sync.RWMutex

	log logger.Logger
}

func NewInMemoryClientConfigRepository(log logger.Logger) repository.ClientConfigRepository {
	return &InMemoryClientConfigRepository{
		storage: make(map[string]*domain.APIClientConfig),
		log:     log,
	}
}

func (r *InMemoryClientConfigRepository) Create(ctx context.Context, config *domain.APIClientConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[config.ID] = config

	r.log.Info("config created",
		"id", config.ID,
		"client_id", config.ClientID,
	)

	return nil
}

func (r *InMemoryClientConfigRepository) GetByID(ctx context.Context, configID string) (*domain.APIClientConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	config, ok := r.storage[configID]
	if !ok {
		r.log.Error("config not found",
			"id", configID,
		)
		return nil, domain.ErrConfigNotFound
	}

	r.log.Info("config fetched",
		"id", configID,
	)

	copyConfig := *config
	return &copyConfig, nil
}

func (r *InMemoryClientConfigRepository) ListByClientID(ctx context.Context, clientID string) ([]*domain.APIClientConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.APIClientConfig, 0, len(r.storage))

	for _, config := range r.storage {
		if config.ClientID == clientID {
			copyConfig := *config
			result = append(result, &copyConfig)
		}
	}

	r.log.Info("list configs by client",
		"client_id", clientID,
		"count", len(result),
	)

	return result, nil
}

func (r *InMemoryClientConfigRepository) Delete(
	ctx context.Context,
	configID string,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.storage[configID]
	if !ok {
		r.log.Error("delete failed: config not found",
			"id", configID,
		)
		return domain.ErrConfigNotFound
	}

	delete(r.storage, configID)

	r.log.Info("config deleted",
		"id", configID,
	)

	return nil
}
