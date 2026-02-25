package repository

import (
	"context"
	"control_plane/internal/domain"
)

type ClientRepository interface {
	Create(ctx context.Context, client *domain.APIClient) error
	GetByID(ctx context.Context, id string) (*domain.APIClient, error)
	Update(ctx context.Context, client *domain.APIClient) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, status string, limit, offset int) ([]*domain.APIClient, int, error)
}