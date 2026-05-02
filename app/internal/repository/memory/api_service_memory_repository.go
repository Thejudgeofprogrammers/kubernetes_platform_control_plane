package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/repository"
	"sync"
)

type InMemoryAPIServiceRepository struct {
	mu      sync.RWMutex
	storage map[string]*domain.APIService

	log logger.Logger
}

func NewInMemoryAPIServiceRepository(log logger.Logger) repository.APIServiceRepository {
	return &InMemoryAPIServiceRepository{
		storage: make(map[string]*domain.APIService),
		log:     log,
	}
}

func (r *InMemoryAPIServiceRepository) Create(ctx context.Context, service *domain.APIService) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[service.ID] = service

	r.log.Info("api service created",
		"id", service.ID,
	)

	return nil
}

func (r *InMemoryAPIServiceRepository) GetByID(ctx context.Context, id string) (*domain.APIService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, ok := r.storage[id]
	if !ok {
		r.log.Error("api service not found",
			"id", id,
		)
		return nil, domain.ErrClientNotFound
	}

	copyService := *service

	r.log.Info("api service fetched",
		"id", id,
	)

	return &copyService, nil
}

func (r *InMemoryAPIServiceRepository) List(ctx context.Context) ([]*domain.APIService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.APIService, 0, len(r.storage))

	for _, s := range r.storage {
		copyService := *s
		result = append(result, &copyService)
	}

	r.log.Info("list api services",
		"count", len(result),
	)

	return result, nil
}

func (r *InMemoryAPIServiceRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[id]; !ok {
		r.log.Error("delete failed: api service not found",
			"id", id,
		)
		return domain.ErrClientNotFound
	}

	delete(r.storage, id)

	r.log.Info("api service deleted",
		"id", id,
	)

	return nil
}

func (r *InMemoryAPIServiceRepository) Update(
	ctx context.Context,
	service *domain.APIService,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.storage[service.ID]
	if !ok {
		return domain.ErrAPIServiceNotFound
	}

	copyService := *service
	r.storage[service.ID] = &copyService

	r.log.Info("api service updated", "id", service.ID)

	return nil
}
