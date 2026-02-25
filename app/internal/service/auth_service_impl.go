package service

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
)

type authService struct {
	userRepo repository.UserRepository
	refreshService RefreshService
}

func NewAuthService() *authService {
	return &authService{}
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	userID, err := s.refreshService.Validate(ctx, refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	_ = s.refreshService.Delete(ctx, refreshToken)

	user, _ := s.userRepo.GetByID(ctx, userID)

	access, _ := s.jwtSer
}
