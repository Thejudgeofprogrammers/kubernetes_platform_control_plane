package service

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/orchestrator"
	"control_plane/internal/repository"
	"time"

	"github.com/google/uuid"
)

type configService struct {
	clientRepo repository.ClientRepository
	configRepo repository.ClientConfigRepository
	audit      *AuditService
	orchestrator orchestrator.Orchestrator
}

func NewConfigService(
	clientRepo repository.ClientRepository,
	configRepo repository.ClientConfigRepository,
	audit *AuditService,
	orchestrator orchestrator.Orchestrator,
) ConfigService {
	return &configService{
		clientRepo: clientRepo,
		configRepo: configRepo,
		audit:      audit,
		orchestrator: orchestrator,
	}
}

func (s *configService) Deploy(
	ctx context.Context,
	userID string,
	clientID string,
	configID string,
) error {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return err
	}

	config, err := s.configRepo.GetByID(ctx, configID)
	if err != nil {
		return err
	}

	if client.ActiveConfigID != nil && *client.ActiveConfigID == configID {
		return nil
	}

	if config.ClientID != clientID {
		return domain.ErrInvalidStateTransition
	}

	if err := s.orchestrator.Deploy(ctx, client, config); err != nil {
		return err
	}

	client.ActivateConfig(configID)

	if err := s.clientRepo.Update(ctx, client); err != nil {
		return err
	}

	return s.audit.Log(ctx, clientID, userID, domain.ActionDeploy)
}

func (s *configService) CreateConfig(
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
) (*domain.APIClientConfig, error) {
	_, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	configs, err := s.configRepo.ListByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	for _, c := range configs {
		if c.Version == version {
			return nil, domain.ErrConfigVersionExists
		}
	}

	config := &domain.APIClientConfig{
		ID:           uuid.NewString(),
		ClientID:     clientID,
		Version:      version,
		AuthType:     authType,
		AuthRef:      authRef,
		TimeoutMs:    timeoutMs,
		RetryCount:   retryCount,
		RetryBackoff: retryBackoff,
		Headers:      headers,
		CreatedAt:    time.Now(),
		CreatedBy:    userID,
	}

	if err := s.configRepo.Create(ctx, config); err != nil {
		return nil, err
	}

	_ = s.audit.Log(ctx, clientID, userID, domain.ActionUpdate)

	return config, nil
}

func (s *configService) ListConfigs(
	ctx context.Context,
	clientID string,
) ([]*domain.APIClientConfig, error) {
	_, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return s.configRepo.ListByClientID(ctx, clientID)
}
