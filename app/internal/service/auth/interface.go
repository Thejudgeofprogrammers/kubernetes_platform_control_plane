package auth

import (
	"context"
	"control_plane/internal/domain"
)

type AuthService interface {
	Register(ctx context.Context, email, fullName string) error
	RequestCode(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, email, code string) (*domain.AuthTokens, error)
	Refresh(ctx context.Context, refreshToken string) (*domain.AuthTokens, error)
}
