package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"control_plane/internal/service/auth"
	"control_plane/internal/service/jwt"
	"control_plane/internal/service/refresh"
)

type authService struct {
	userRepo repository.UserRepository
	refreshService refresh.RefreshService
	jwtService jwt.JWTService
}

func NewAuthService(userRepo repository.UserRepository, refreshService refresh.RefreshService, jwtService jwt.JWTService) auth.AuthService {
	return &authService{userRepo: userRepo, refreshService: refreshService, jwtService: jwtService}
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*jwt.AuthTokens, error) {
	userID, err := s.refreshService.Validate(ctx, refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidRefreshToken
	}

	if err := s.refreshService.Delete(ctx, refreshToken); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	access, err := s.jwtService.GenerateAccessToken(user.ID, string(user.Role))
	if err != nil {
		return nil, err
	}
	
	newRefresh, err := s.refreshService.Create(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &jwt.AuthTokens{
		AccessToken: access,
		RefreshToken: newRefresh,
	}, nil
}

func (s *authService) Register(ctx context.Context, email, fullName string) error {

}

func (s *authService) RequestCode(ctx context.Context, email string) error {

}

func (s *authService) VerifyCode(ctx context.Context, email, code string) (*jwt.AuthTokens, error) {

}
