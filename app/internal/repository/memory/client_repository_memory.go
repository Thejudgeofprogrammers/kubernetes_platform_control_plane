package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/repository"
	"sort"
	"sync"
)

type InMemoryClientRepository struct {
	storage map[string]*domain.APIClient
	mu      sync.RWMutex

	log logger.Logger
}

func NewInMemoryClientRepository(log logger.Logger) repository.ClientRepository {
	return &InMemoryClientRepository{
		storage: make(map[string]*domain.APIClient),
		log:     log,
	}
}

func (r *InMemoryClientRepository) Create(ctx context.Context, client *domain.APIClient) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.storage[client.ID]; exists {
		r.log.Error("client already exists", "id", client.ID)
		return domain.ErrClientAlredyExist
	}

	r.storage[client.ID] = client
	r.log.Info("client created", "id", client.ID)
	return nil
}

func (r *InMemoryClientRepository) GetByID(ctx context.Context, id string) (*domain.APIClient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client, ok := r.storage[id]
	if !ok {
		r.log.Error("client not found", "id", id)
		return nil, domain.ErrClientNotFound
	}

	r.log.Info("client fetched", "id", id)
	return client, nil
}

func (r *InMemoryClientRepository) Update(ctx context.Context, client *domain.APIClient) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[client.ID]; !ok {
		r.log.Error("update failed: client not found", "id", client.ID)
		return domain.ErrClientNotFound
	}

	r.storage[client.ID] = client
	r.log.Info("client updated", "id", client.ID)
	return nil
}

func (r *InMemoryClientRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.storage[id]; !ok {
		r.log.Error("delete failed: client not found", "id", id)
		return domain.ErrClientNotFound
	}
	delete(r.storage, id)
	r.log.Info("client deleted", "id", id)
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
		copyClient := *client
		result = append(result, &copyClient)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	total := len(result)

	if offset >= len(result) {
		r.log.Info("list clients: empty result due to offset", "offset", offset, "total", total)
		return []*domain.APIClient{}, total, nil
	}

	result = result[offset:]

	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}

	r.log.Info("list clients",
		"status", status,
		"limit", limit,
		"offset", offset,
		"returned", len(result),
		"total", total,
	)

	return result, total, nil
}
