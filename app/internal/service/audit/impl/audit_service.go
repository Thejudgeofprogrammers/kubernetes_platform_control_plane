package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"control_plane/internal/service/audit"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type auditService struct {
	repo repository.ClientActionRepository
	log  *slog.Logger
}

func NewAuditService(repo repository.ClientActionRepository, log *slog.Logger) audit.AuditService {
	return &auditService{
		repo: repo,
		log:  log,
	}
}

func (s *auditService) Log(
	ctx context.Context,
	clientID string,
	userID string,
	action domain.ActionType,
) error {
	record := &domain.APIClientAction{
		ID:        uuid.NewString(),
		ClientID:  clientID,
		UserID:    userID,
		Type:      action,
		CreatedAt: time.Now(),
	}

	return s.repo.Create(ctx, record)
}
