package apiservice

import (
	"context"
	"control_plane/internal/domain"
)

type APIServiceService interface {
	Create(ctx context.Context, name, baseURL, protocol string) (*domain.APIService, error)
	List(ctx context.Context) ([]*domain.APIService, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*domain.APIService, error)
	Update(ctx context.Context, id, name, baseURL, protocol string) (*domain.APIService, error)
}
