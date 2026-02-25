package memory

import (
	"context"
	"control_plane/internal/domain"
	"sync"
)

type InMemoryClientActionRepository struct {
	storage []*domain.APIClientAction
	mu sync.RWMutex
}

func NewInMemoryClientActionRepository() *InMemoryClientActionRepository {
	return &InMemoryClientActionRepository{
		storage: []*domain.APIClientAction{},
	}
}

func (r *InMemoryClientActionRepository) Create(ctx context.Context, action *domain.APIClientAction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage = append(r.storage, action)
	return nil
}

func (r *InMemoryClientActionRepository) ListByClientID(ctx context.Context, clientID string) ([]*domain.APIClientAction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.APIClientAction, 0, len(r.storage))

	for _, action := range r.storage {
		if action.ClientID == clientID {
			result = append(result, action)
		}
	}

	return result, nil
}