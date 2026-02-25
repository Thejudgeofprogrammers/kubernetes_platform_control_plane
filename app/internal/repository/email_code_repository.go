package repository

import (
	"context"
	"control_plane/internal/domain"
)

type EmailCodeRepository interface {
	Save(ctx context.Context, code *domain.EmailCode) error
	Get(ctx context.Context, email string) (*domain.EmailCode, error)
	Delete(ctx context.Context, email string) error
}
