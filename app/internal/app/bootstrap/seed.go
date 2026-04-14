package bootstrap

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

func SeedAdmin(
	ctx context.Context,
	userRepo repository.UserRepository,
	email, fullname string,
) error {

	users, err := userRepo.List(ctx)
	if err != nil {
		slog.Error("failed to list users", "error", err)
		return err
	}

	if len(users) > 0 {
		slog.Debug("admin seed skipped", "users_count", len(users))
		return nil
	}

	admin := &domain.User{
		ID:        uuid.NewString(),
		Email:     email,
		FullName:  fullname,
		Role:      domain.RoleOwner,
		CreatedAt: time.Now(),
	}

	if err := userRepo.Create(ctx, admin); err != nil {
		slog.Error("failed to create admin", "error", err)
		return err
	}

	slog.Info("admin user created",
		"email", email,
	)

	return nil
}