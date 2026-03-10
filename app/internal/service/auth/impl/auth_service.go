package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"control_plane/internal/service/auth"
	"control_plane/internal/service/jwt"
	"control_plane/internal/service/refresh"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

type authService struct {
	userRepo       repository.UserRepository
	codeRepo       repository.EmailCodeRepository
	refreshService refresh.RefreshService
	jwtService     jwt.JWTService
	expire         int
}

func NewAuthService(
	userRepo repository.UserRepository,
	codeRepo repository.EmailCodeRepository,
	refreshService refresh.RefreshService,
	jwtService jwt.JWTService,
	exp int,
) auth.AuthService {
	return &authService{
		userRepo:       userRepo,
		codeRepo:       codeRepo,
		refreshService: refreshService,
		jwtService:     jwtService,
		expire: exp,
	}
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*domain.AuthTokens, error) {
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

	access, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	newRefresh, err := s.refreshService.Create(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &domain.AuthTokens{
		AccessToken:  access,
		RefreshToken: newRefresh,
	}, nil
}

func (s *authService) Register(ctx context.Context, email, fullName string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && user != nil {
		return domain.ErrUserAlreadyExists
	}

	if err != nil && err != domain.ErrUserNotFound {
		return err
	}

	newUser := &domain.User{
		ID:        uuid.NewString(),
		Email:     email,
		FullName:  fullName,
		Role:      domain.RoleViewer,
		CreatedAt: time.Now(),
	}

	return s.userRepo.Create(ctx, newUser)
}

func (s *authService) RequestCode(ctx context.Context, email string) error {
	email = strings.ToLower(email)
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user == nil {
		return domain.ErrUserNotFound
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	emailCode := &domain.EmailCode{
		Email: email,
		Code: code,
		ExpiresAt: time.Now().Add(time.Duration(s.expire) * time.Second),
	}

	return s.codeRepo.Save(ctx, emailCode)
}

func (s *authService) VerifyCode(ctx context.Context, email, code string) (*domain.AuthTokens, error) {
	email = strings.ToLower(email)
	storedCode, err := s.codeRepo.Get(ctx, email)
	if err != nil {
		return nil, err
	}

	if storedCode.Code != code {
		return nil, domain.ErrCodeNotFound
	}

	if err := s.codeRepo.Delete(ctx, email); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	access, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refresh, err := s.refreshService.Create(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &domain.AuthTokens{
		AccessToken: access,
		RefreshToken: refresh,
	}, nil
}
