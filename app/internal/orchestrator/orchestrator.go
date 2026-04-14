package orchestrator

import (
	"context"
	"control_plane/internal/domain"
)

type Orchestrator interface {
	Deploy(ctx context.Context, client *domain.APIClient, config *domain.APIClientConfig) error
	Restart(ctx context.Context, clientID string) error
	Delete(ctx context.Context, clientID string) error
	CheckHealth(ctx context.Context, clientID string)
}