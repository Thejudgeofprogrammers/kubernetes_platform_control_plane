package audit

import (
	"context"
	"control_plane/internal/domain"
)

type AuditService interface {
	Log(ctx context.Context, clientID string,userID string, action domain.ActionType) error 
}
