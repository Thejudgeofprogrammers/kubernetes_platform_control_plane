package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/repository"
	"control_plane/internal/service/action"
	"time"

	"github.com/google/uuid"
)

type actionService struct {
	repo repository.ClientActionRepository
	log  logger.Logger
}

func NewActionService(repo repository.ClientActionRepository, log logger.Logger) action.ActionService {
	return &actionService{
		repo: repo,
		log:  log,
	}
}

func (s *actionService) Create(ctx context.Context, clientID, userID string, actionType domain.ActionType) (*domain.APIClientAction, error) {

	s.log.Info("action create started",
		"client_id", clientID,
		"user_id", userID,
		"type", actionType,
	)

	action := &domain.APIClientAction{
		ID:        uuid.NewString(),
		ClientID:  clientID,
		UserID:    userID,
		Type:      actionType,
		Status:    domain.ActionPending,
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(ctx, action)
	if err != nil {
		s.log.Error("failed to create action",
			"action_id", action.ID,
			"client_id", clientID,
			"user_id", userID,
			"type", actionType,
			"error", err,
		)
		return nil, err
	}

	s.log.Info("action created",
		"action_id", action.ID,
		"client_id", clientID,
		"user_id", userID,
		"type", actionType,
	)

	return action, nil
}
