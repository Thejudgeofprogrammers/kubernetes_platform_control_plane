package service

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"time"

	"github.com/google/uuid"
)

type AuditService struct {
	repo repository.ClientActionRepository
}

func NewAuditService(repo repository.ClientActionRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) Log(
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
