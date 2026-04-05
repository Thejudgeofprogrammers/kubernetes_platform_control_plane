package repository

import (
	"context"
	"control_plane/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)

	List(ctx context.Context) ([]domain.User, error)
	Delete(ctx context.Context, id string) error
	UpdateRole(ctx context.Context, id string, role string) error
}
