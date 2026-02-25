package memory

import (
	"context"
	"control_plane/internal/domain"
	"sort"
	"sync"
)

type InMemoryClientRepository struct {
	storage map[string]*domain.APIClient
	mu sync.RWMutex
}

func NewInMemoryClientRepository() *InMemoryClientRepository {
	return &InMemoryClientRepository{
		storage: make(map[string]*domain.APIClient),
	}
}

func (r *InMemoryClientRepository) Create(ctx context.Context, client *domain.APIClient) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[client.ID] = client
	return nil
}

func (r *InMemoryClientRepository) GetByID(ctx context.Context, id string) (*domain.APIClient, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	client, ok := r.storage[id]
	if !ok {
		return nil, domain.ErrClientNotFound
	}

	return client, nil
}

func (r *InMemoryClientRepository) Update(ctx context.Context, client *domain.APIClient) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.storage[client.ID] = client
	return nil
}

func (r *InMemoryClientRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, ok := r.storage[id]; !ok {
		return domain.ErrClientNotFound
	}
	delete(r.storage, id)
	return nil
}

func (r *InMemoryClientRepository) List(ctx context.Context, status string, limit, offset int) ([]*domain.APIClient, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	

	result := make([]*domain.APIClient, 0, len(r.storage))

	for _, client := range r.storage {
		if status != "" && string(client.GetStatus()) != status {
			continue
		}
		result = append(result, client)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	total := len(result)

	if offset >= len(result) {
		return []*domain.APIClient{}, total, nil
	}

	result = result[offset:]

	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}

	return result, total, nil
}