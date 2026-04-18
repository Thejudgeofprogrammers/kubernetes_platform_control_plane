package repository

import (
	"context"
	"control_plane/internal/domain"
)

type ClientConfigRepository interface {
	Create(ctx context.Context, config *domain.APIClientConfig) error
	GetByID(ctx context.Context, configID string) (*domain.APIClientConfig, error)
	ListByClientID(ctx context.Context, clientID string) ([]*domain.APIClientConfig, error)
	Delete(ctx context.Context, configID string) error
}