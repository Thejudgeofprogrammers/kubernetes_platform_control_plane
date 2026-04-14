package repository

import (
	"context"
	"control_plane/internal/domain"
)

type ClientHealthRepostiory interface {
	Set(clientID string, health domain.APIClientHealth)
	Update(ctx context.Context, health *domain.APIClientHealth) error 
	GetByClientID(ctx context.Context, clientID string) (*domain.APIClientHealth, error)
}
