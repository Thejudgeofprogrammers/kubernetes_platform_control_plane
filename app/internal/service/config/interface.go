package config

import (
	"context"
	"control_plane/internal/domain"
	configDTO "control_plane/internal/transport/http_gin/dto/config"
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

	Delete(
		ctx context.Context,
		clientID string,
		configID string,
	) error

	Update(
		ctx context.Context,
		userID string,
		clientID string,
		configID string,
		req configDTO.ClientConfigRequest,
	) (*domain.APIClientConfig, error)
}
