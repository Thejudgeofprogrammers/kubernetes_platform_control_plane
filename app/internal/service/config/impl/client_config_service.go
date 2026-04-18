package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/orchestrator"
	"control_plane/internal/repository"
	configDTO "control_plane/internal/transport/http_gin/dto/config"
	"control_plane/internal/service/audit"
	"control_plane/internal/service/config"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type configService struct {
	clientRepo   repository.ClientRepository
	configRepo   repository.ClientConfigRepository
	actionRepo   repository.ClientActionRepository
	audit        audit.AuditService
	orchestrator orchestrator.Orchestrator
	log          *slog.Logger
}

func NewConfigService(
	clientRepo repository.ClientRepository,
	configRepo repository.ClientConfigRepository,
	audit audit.AuditService,
	actionRepo   repository.ClientActionRepository,
	orchestrator orchestrator.Orchestrator,
	log *slog.Logger,
) config.ConfigService {
	return &configService{
		clientRepo:   clientRepo,
		configRepo:   configRepo,
		audit:        audit,
		actionRepo:   actionRepo,
		orchestrator: orchestrator,
		log:          log,
	}
}

func (s *configService) Deploy(
	ctx context.Context,
	userID string,
	clientID string,
	configID string,
) error {
	s.log.Info("deploy config started",
		"user_id", userID,
		"client_id", clientID,
		"config_id", configID,
	)

	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		s.log.Error("failed to get client",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	config, err := s.configRepo.GetByID(ctx, configID)
	if err != nil {
		s.log.Error("failed to get config",
			"client_id", clientID,
			"config_id", configID,
			"error", err,
		)
		return err
	}

	if client.ActiveConfigID != nil && *client.ActiveConfigID == configID {
		s.log.Info("config already active, skipping deploy",
			"client_id", clientID,
			"config_id", configID,
		)
		return nil
	}

	if config.ClientID != clientID {
		s.log.Error("config does not belong to client",
			"client_id", clientID,
			"config_id", configID,
		)
		return domain.ErrInvalidStateTransition
	}

	s.log.Info("orchestrator deploy",
		"client_id", clientID,
		"config_id", configID,
	)

	if err := s.orchestrator.Deploy(ctx, client, config); err != nil {
		s.log.Error("deploy failed",
			"client_id", clientID,
			"config_id", configID,
			"error", err,
		)
		return err
	}

	client.ActivateConfig(configID)

	if err := s.clientRepo.Update(ctx, client); err != nil {
		s.log.Error("failed to update client after deploy",
			"client_id", clientID,
			"config_id", configID,
			"error", err,
		)
		return err
	}

	if err := s.audit.Log(ctx, clientID, userID, domain.ActionDeploy); err != nil {
		s.log.Error("failed to write audit log",
			"client_id", clientID,
			"user_id", userID,
			"error", err,
		)
	}

	s.log.Info("deploy config completed",
		"user_id", userID,
		"client_id", clientID,
		"config_id", configID,
	)

	return nil
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
	s.log.Info("create config started",
		"user_id", userID,
		"client_id", clientID,
		"version", version,
	)

	_, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		s.log.Error("failed to get client",
			"client_id", clientID,
			"error", err,
		)
		return nil, err
	}

	configs, err := s.configRepo.ListByClientID(ctx, clientID)
	if err != nil {
		s.log.Error("failed to list configs",
			"client_id", clientID,
			"error", err,
		)
		return nil, err
	}

	for _, c := range configs {
		if c.Version == version {
			s.log.Warn("config version already exists",
				"client_id", clientID,
				"version", version,
			)
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
		s.log.Error("failed to create config",
			"client_id", clientID,
			"config_id", config.ID,
			"error", err,
		)
		return nil, err
	}

	if err := s.audit.Log(ctx, clientID, userID, domain.ActionUpdate); err != nil {
		s.log.Error("failed to write audit log",
			"client_id", clientID,
			"user_id", userID,
			"error", err,
		)
	}

	s.log.Info("config created",
		"client_id", clientID,
		"config_id", config.ID,
		"version", version,
	)

	return config, nil
}

func (s *configService) ListConfigs(
	ctx context.Context,
	clientID string,
) ([]*domain.APIClientConfig, error) {

	s.log.Debug("list configs",
		"client_id", clientID,
	)

	_, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		s.log.Error("failed to get client",
			"client_id", clientID,
			"error", err,
		)
		return nil, err
	}

	configs, err := s.configRepo.ListByClientID(ctx, clientID)
	if err != nil {
		s.log.Error("failed to list configs",
			"client_id", clientID,
			"error", err,
		)
		return nil, err
	}

	s.log.Debug("configs listed",
		"client_id", clientID,
		"count", len(configs),
	)

	return configs, nil
}

func (s *configService) Delete(
	ctx context.Context,
	clientID string,
	configID string,
) error {

	s.log.Info("delete config started",
		"client_id", clientID,
		"config_id", configID,
	)

	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return err
	}

	config, err := s.configRepo.GetByID(ctx, configID)
	if err != nil {
		return err
	}

	if config.ClientID != client.ID {
		return domain.ErrConfigNotFound
	}

	if client.ActiveConfigID != nil && *client.ActiveConfigID == configID {
		return domain.ErrInvalidStateTransition
	}

	if err := s.configRepo.Delete(ctx, configID); err != nil {
		return err
	}

	s.log.Info("config deleted",
		"client_id", clientID,
		"config_id", configID,
	)

	return nil
}

func (s *configService) Update(
	ctx context.Context,
	userID string,
	clientID string,
	configID string,
	req configDTO.ClientConfigRequest,
) (*domain.APIClientConfig, error) {

	s.log.Info("update config started",
		"client_id", clientID,
		"config_id", configID,
		"user_id", userID,
	)

	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	oldConfig, err := s.configRepo.GetByID(ctx, configID)
	if err != nil {
		return nil, err
	}

	if oldConfig.ClientID != client.ID {
		return nil, domain.ErrConfigNotFound
	}

	newConfig := &domain.APIClientConfig{
		ID:           uuid.NewString(),
		ClientID:     clientID,
		Version:      req.Version,
		AuthType:     domain.AuthType(req.AuthType),
		AuthRef:      req.AuthRef,
		TimeoutMs:    req.TimeoutMs,
		RetryCount:   req.RetryCount,
		RetryBackoff: req.RetryBackoff,
		Headers:      req.Headers,
		CreatedAt:    time.Now(),
		CreatedBy:    userID,
	}

	if err := s.configRepo.Create(ctx, newConfig); err != nil {
		return nil, err
	}

	if client.ActiveConfigID != nil && *client.ActiveConfigID == configID {

		client.ActivateConfig(newConfig.ID)

		if err := s.clientRepo.Update(ctx, client); err != nil {
			return nil, err
		}

		action := &domain.APIClientAction{
			ID:        uuid.NewString(),
			ClientID:  clientID,
			UserID:    userID,
			Type:      domain.ActionDeploy,
			CreatedAt: time.Now(),
		}

		if err := s.actionRepo.Create(ctx, action); err != nil {
			return nil, err
		}

		s.log.Info("deploy triggered after config update",
			"client_id", clientID,
			"new_config_id", newConfig.ID,
		)
	}

	s.log.Info("config updated (new version created)",
		"client_id", clientID,
		"old_config_id", configID,
		"new_config_id", newConfig.ID,
	)

	return newConfig, nil
}
