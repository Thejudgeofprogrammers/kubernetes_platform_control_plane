package repository

import (
	"context"
	"control_plane/internal/domain"
)

type ClientHealthRepostiory interface {
	Update(ctx context.Context, health *domain.APIClientHealth) error 
	GetByClientID(ctx context.Context, clientID string) (*domain.APIClientHealth, error)
}
