package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"log/slog"
	"sync"
)

type InMemoryClientHealthRepostiory struct {
	mu      sync.RWMutex
	storage map[string]*domain.APIClientHealth

	log *slog.Logger
}

func NewInMemoryClientHealthRepository(log *slog.Logger) repository.ClientHealthRepostiory {
	return &InMemoryClientHealthRepostiory{
		storage: make(map[string]*domain.APIClientHealth),
		log:     log,
	}
}

func (r *InMemoryClientHealthRepostiory) Update(ctx context.Context, health *domain.APIClientHealth) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[health.ClientID] = health

	r.log.Info("client health updated",
		"client_id", health.ClientID,
	)

	return nil
}

func (r *InMemoryClientHealthRepostiory) GetByClientID(ctx context.Context, clientID string) (*domain.APIClientHealth, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	health, ok := r.storage[clientID]
	if !ok {
		r.log.Info("client health not found",
			"client_id", clientID,
		)
		return nil, nil
	}

	copyHealth := *health

	r.log.Info("client health fetched",
		"client_id", clientID,
	)

	return &copyHealth, nil
}

func (r *InMemoryClientHealthRepostiory) Set(clientID string, health domain.APIClientHealth) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.storage[clientID] = &health

	r.log.Info("client health updated",
		"client_id", health.ClientID,
	)
}