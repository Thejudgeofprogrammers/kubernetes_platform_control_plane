package repository

import (
    "context"
    "control_plane/internal/domain"
)

type APIServiceRepository interface {
    Create(ctx context.Context, service *domain.APIService) error
    GetByID(ctx context.Context, id string) (*domain.APIService, error)
    List(ctx context.Context) ([]*domain.APIService, error)
    Delete(ctx context.Context, id string) error
    Update(ctx context.Context, service *domain.APIService) error
}