package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/repository"
	"control_plane/internal/service/auth"
	"control_plane/internal/service/email"
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
	emailSender    email.EmailSender
	expire         int
	log            logger.Logger
}

func NewAuthService(
	userRepo repository.UserRepository,
	codeRepo repository.EmailCodeRepository,
	refreshService refresh.RefreshService,
	jwtService jwt.JWTService,
	emailSender email.EmailSender,
	exp int,
	log logger.Logger,
) auth.AuthService {
	return &authService{
		userRepo:       userRepo,
		codeRepo:       codeRepo,
		refreshService: refreshService,
		jwtService:     jwtService,
		emailSender:    emailSender,
		expire:         exp,
		log:            log,
	}
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*domain.AuthTokens, error) {
	s.log.Info("refresh flow started")

	userID, err := s.refreshService.Validate(ctx, refreshToken)
	if err != nil {
		s.log.Warn("invalid refresh token")
		return nil, domain.ErrInvalidRefreshToken
	}

	if err := s.refreshService.Delete(ctx, refreshToken); err != nil {
		s.log.Error("failed to delete refresh token",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}

	access, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		s.log.Error("failed to generate access token",
			"user_id", user.ID,
			"error", err,
		)
		return nil, err
	}

	newRefresh, err := s.refreshService.Create(ctx, user.ID)
	if err != nil {
		s.log.Error("failed to create refresh token",
			"user_id", user.ID,
			"error", err,
		)
		return nil, err
	}

	s.log.Info("refresh flow completed",
		"user_id", user.ID,
	)

	return &domain.AuthTokens{
		AccessToken:  access,
		RefreshToken: newRefresh,
	}, nil
}

func (s *authService) Register(ctx context.Context, email, fullName string) error {
	s.log.Info("register user started",
		"email", email,
	)

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && user != nil {
		s.log.Warn("user already exists",
			"email", email,
		)
		return domain.ErrUserAlreadyExists
	}

	if err != nil && err != domain.ErrUserNotFound {
		s.log.Error("failed to check user existence",
			"email", email,
			"error", err,
		)
		return err
	}

	newUser := &domain.User{
		ID:        uuid.NewString(),
		Email:     email,
		FullName:  fullName,
		Role:      domain.RoleViewer,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		s.log.Error("failed to create user",
			"email", email,
			"error", err,
		)
		return err
	}

	s.log.Info("user registered",
		"user_id", newUser.ID,
		"email", email,
	)

	return nil
}

func (s *authService) RequestCode(ctx context.Context, email string) error {
	email = strings.ToLower(email)

	s.log.Info("request code started",
		"email", email,
	)

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user",
			"email", email,
			"error", err,
		)
		return err
	}

	if user == nil {
		s.log.Warn("user not found for code request",
			"email", email,
		)
		return domain.ErrUserNotFound
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	emailCode := &domain.EmailCode{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(time.Duration(s.expire) * time.Second),
	}

	if err := s.codeRepo.Save(ctx, emailCode); err != nil {
		s.log.Error("failed to save email code",
			"email", email,
			"error", err,
		)
		return err
	}

	if err := s.emailSender.Send(email, code); err != nil {
		s.log.Error("failed to send email",
			"email", email,
			"error", err,
		)
		return err
	}

	return nil
}

func (s *authService) VerifyCode(ctx context.Context, email, code string) (*domain.AuthTokens, error) {
	email = strings.ToLower(email)

	s.log.Info("verify code started",
		"email", email,
	)

	storedCode, err := s.codeRepo.Get(ctx, email)
	if err != nil {
		s.log.Error("failed to get stored code",
			"email", email,
			"error", err,
		)
		return nil, err
	}

	if storedCode.Code != code {
		s.log.Warn("invalid verification code",
			"email", email,
		)
		return nil, domain.ErrCodeNotFound
	}

	if err := s.codeRepo.Delete(ctx, email); err != nil {
		s.log.Error("failed to delete code",
			"email", email,
			"error", err,
		)
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.log.Error("failed to get user",
			"email", email,
			"error", err,
		)
		return nil, err
	}

	access, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		s.log.Error("failed to generate access token",
			"user_id", user.ID,
			"error", err,
		)
		return nil, err
	}

	refresh, err := s.refreshService.Create(ctx, user.ID)
	if err != nil {
		s.log.Error("failed to create refresh token",
			"user_id", user.ID,
			"error", err,
		)
		return nil, err
	}

	s.log.Info("verify code completed",
		"user_id", user.ID,
		"email", email,
	)

	return &domain.AuthTokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
