package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"sync"
)

type InMemoryAPIServiceRepository struct {
	mu sync.RWMutex
	storage map[string]*domain.APIService
}

func NewInMemoryAPIServiceRepository() repository.APIServiceRepository {
	return &InMemoryAPIServiceRepository{
		storage: make(map[string]*domain.APIService),
	}
}

func (r *InMemoryAPIServiceRepository) Create(ctx context.Context, service *domain.APIService) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[service.ID] = service
	return nil
}

func (r *InMemoryAPIServiceRepository) GetByID(ctx context.Context, id string) (*domain.APIService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, ok := r.storage[id]
	if !ok {
		return nil, domain.ErrClientNotFound
	}

	return service, nil
}

func (r *InMemoryAPIServiceRepository) List(ctx context.Context) ([]*domain.APIService, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.APIService, 0, len(r.storage))

	for _, s := range r.storage {
		result = append(result, s)
	}

	return result, nil
}

func (r *InMemoryAPIServiceRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.storage, id)
	return nil
}