package repository

import (
	"context"
	"control_plane/internal/domain"
)

type ClientAccessRepository interface {
	Grant(ctx context.Context, access *domain.AuthClientAccess) error
	ListByUserID(ctx context.Context, userID string) ([]domain.AuthClientAccess, error)
	Get(ctx context.Context, userID, clientID string) (*domain.AuthClientAccess, error)
}
