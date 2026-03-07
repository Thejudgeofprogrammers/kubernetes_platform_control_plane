package auth

import (
	"context"
	"control_plane/internal/service/jwt"
)

type AuthService interface {
	Register(ctx context.Context, email, fullName string) error
	RequestCode(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, email, code string) (*jwt.AuthTokens, error)
	Refresh(ctx context.Context, refreshToken string) (*jwt.AuthTokens, error)
}
