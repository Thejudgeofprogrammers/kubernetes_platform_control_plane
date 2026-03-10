package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"sync"
)

type InMemoryClientHealthRepostiory struct {
	mu      sync.RWMutex
	storage map[string]*domain.APIClientHealth
}

func NewInMemoryClientHealthRepository() repository.ClientHealthRepostiory {
	return &InMemoryClientHealthRepostiory{
		storage: make(map[string]*domain.APIClientHealth),
	}
}

func (r *InMemoryClientHealthRepostiory) Update(ctx context.Context, health *domain.APIClientHealth) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[health.ClientID] = health
	return nil
}

func (r *InMemoryClientHealthRepostiory) GetByClientID(ctx context.Context, clientID string) (*domain.APIClientHealth, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	health, ok := r.storage[clientID]
	if !ok {
		return nil, nil
	}

	return health, nil
}
