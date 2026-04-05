package action

import (
	"context"
	"control_plane/internal/domain"
)

type ActionService interface {
	Create(ctx context.Context, clientID, userID string, actionType domain.ActionType) (*domain.APIClientAction, error)
}
