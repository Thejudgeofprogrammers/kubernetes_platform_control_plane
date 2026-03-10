package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"sort"
	"sync"
)

type InMemoryClientActionRepository struct {
	storage []*domain.APIClientAction
	mu sync.RWMutex
}

func NewInMemoryClientActionRepository() repository.ClientActionRepository {
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
			copyAction := *action
			result = append(result, &copyAction)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result, nil
}