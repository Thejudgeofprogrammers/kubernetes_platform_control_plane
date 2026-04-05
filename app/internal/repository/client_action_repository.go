package repository

import (
	"context"
	"control_plane/internal/domain"
)

type ClientActionRepository interface {
	Create(ctx context.Context, action *domain.APIClientAction) error
	ListByClientID(ctx context.Context, clientID string) ([]*domain.APIClientAction, error)
	GetPending(ctx context.Context) ([]*domain.APIClientAction, error)
	UpdateStatus(ctx context.Context, id string, status domain.ActionStatus, err *string) error
}
