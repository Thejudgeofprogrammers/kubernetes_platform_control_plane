package service

import "context"

type AuthService interface {
	Register(ctx context.Context, email, fullName string) error
	RequestCode(ctx context.Context, email string) error
	VerifyCode(ctx context.Context, email, code string) (*AuthTokens, error)
	Refresh(ctx context.Context, refreshToken string) (*AuthTokens, error)
}
