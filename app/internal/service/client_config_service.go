package service

import (
	"context"
	"control_plane/internal/domain"
)

type ConfigService interface {
	CreateConfig(
		ctx context.Context,
		userID string,
		clientID string,
		version string,
		authType domain.AuthType,
		authRef string,
		timeoutMs int,
		retryCount int,
		retryBackoff int,
		headers map[string]string,
	) (*domain.APIClientConfig, error)

	ListConfigs(
		ctx context.Context,
		clientID string,
	) ([]*domain.APIClientConfig, error)

	Deploy(
		ctx context.Context,
		userID string,
		clientID string,
		configID string,
	) error
}
