package client

import (
	"context"
	"control_plane/internal/domain"
)

type ClientService interface {
	List(ctx context.Context, status string, limit, offset int) ([]domain.APIClient, int, error)
	Create(ctx context.Context, userID, name, apiServiceID, description string) (*domain.APIClient, error)
	GetByID(ctx context.Context, id string) (*domain.APIClient, error)
	Restart(ctx context.Context, userID, id string, reason string) error
	Delete(ctx context.Context, userID, id string) error
	Start(ctx context.Context, clientID string) error
}