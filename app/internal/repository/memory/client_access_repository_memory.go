package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"log/slog"
	"sync"
)

// Resource role not implement

type InMemoryClientAccessRepository struct {
	mu      sync.RWMutex
	storage []*domain.AuthClientAccess

	log *slog.Logger
}

func NewInMemoryClientAccessRepository(log *slog.Logger) repository.ClientAccessRepository {
	return &InMemoryClientAccessRepository{
		log: log,
	}
}

func (r *InMemoryClientAccessRepository) Grant(ctx context.Context, access *domain.AuthClientAccess) error {
	return nil
}

func (r *InMemoryClientAccessRepository) ListByUserID(ctx context.Context, userID string) ([]domain.AuthClientAccess, error) {
	return nil, nil
}

func (r *InMemoryClientAccessRepository) Get(ctx context.Context, userID, clientID string) (*domain.AuthClientAccess, error) {
	return nil, nil
}
