package user

import (
	"context"
	"control_plane/internal/domain"
)

type UserService interface {
	List(ctx context.Context) ([]domain.User, error)
	Delete(ctx context.Context, id string) error
	UpdateRole(ctx context.Context, id string, role string) error
	GetMe(ctx context.Context, userID string) (*domain.User, error)
}
