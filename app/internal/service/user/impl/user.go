package impl

import (
	"context"

	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/repository"
	"control_plane/internal/service/user"
)

type userService struct {
	userRepo repository.UserRepository
	log      logger.Logger
}

func NewUserService(userRepo repository.UserRepository, log logger.Logger) user.UserService {
	return &userService{
		userRepo: userRepo,
		log:      log,
	}
}

func (s *userService) List(ctx context.Context) ([]domain.User, error) {
	s.log.Debug("list users started")

	users, err := s.userRepo.List(ctx)
	if err != nil {
		s.log.Error("failed to list users",
			"error", err,
		)
		return nil, err
	}

	s.log.Debug("users listed",
		"count", len(users),
	)

	return users, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	requestUserID, _ := ctx.Value("user_id").(string)

	s.log.Info("delete user started",
		"target_user_id", id,
		"requested_by", requestUserID,
	)

	if id == "" {
		s.log.Warn("empty user id on delete")
		return domain.ErrEmptyUserID
	}

	// Проверка на удаления самого себя:
	// if id == requestUserID { ... }

	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete user",
			"target_user_id", id,
			"error", err,
		)
		return err
	}

	s.log.Info("user deleted",
		"target_user_id", id,
		"requested_by", requestUserID,
	)

	return nil
}

func (s *userService) UpdateRole(ctx context.Context, id string, role string) error {
	requestUserID, _ := ctx.Value("user_id").(string)

	s.log.Info("update user role started",
		"target_user_id", id,
		"new_role", role,
		"requested_by", requestUserID,
	)

	if id == "" {
		s.log.Warn("empty user id on role update")
		return domain.ErrEmptyUserID
	}

	switch role {
	case string(domain.RoleOwner),
		string(domain.RoleOperator),
		string(domain.RoleViewer):
	default:
		s.log.Warn("invalid role",
			"role", role,
			"target_user_id", id,
		)
		return domain.ErrInvalidRole
	}

	if err := s.userRepo.UpdateRole(ctx, id, role); err != nil {
		s.log.Error("failed to update role",
			"target_user_id", id,
			"new_role", role,
			"error", err,
		)
		return err
	}

	s.log.Info("user role updated",
		"target_user_id", id,
		"new_role", role,
		"requested_by", requestUserID,
	)

	return nil
}

func (s *userService) GetMe(ctx context.Context, userID string) (*domain.User, error) {
	s.log.Debug("get current user",
		"user_id", userID,
	)

	if userID == "" {
		s.log.Warn("unauthorized get me")
		return nil, domain.ErrUnAuthorizedUser
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}

	return user, nil
}
