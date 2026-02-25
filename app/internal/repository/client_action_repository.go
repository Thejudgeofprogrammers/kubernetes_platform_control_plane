package repository

import (
	"context"
	"control_plane/internal/domain"
)

type ClientActionRepository interface {
	Create(ctx context.Context, action *domain.APIClientAction) error
	ListByClientID(ctx context.Context, clientID string) ([]*domain.APIClientAction, error)
}
