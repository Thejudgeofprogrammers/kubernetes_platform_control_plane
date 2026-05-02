package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/orchestrator"
	"control_plane/internal/repository"
	"control_plane/internal/service/config"
	configDTO "control_plane/internal/transport/http_gin/dto/config"
	actionSrv "control_plane/internal/service/action"
	"time"

	"github.com/google/uuid"
)

type configService struct {
	clientRepo   repository.ClientRepository
	configRepo   repository.ClientConfigRepository
	actionService    actionSrv.ActionService
	orchestrator orchestrator.Orchestrator
	log          logger.Logger
}

func NewConfigService(
	clientRepo repository.ClientRepository,
	configRepo repository.ClientConfigRepository,
	actionService actionSrv.ActionService,
	orchestrator orchestrator.Orchestrator,
	log logger.Logger,
) config.ConfigService {
	return &configService{
		clientRepo:   clientRepo,
		configRepo:   configRepo,
		actionService: actionService,
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

	if config.ClientID != clientID {
		s.log.Error("config does not belong to client",
			"client_id", clientID,
			"config_id", configID,
		)
		return domain.ErrInvalidStateTransition
	}

	actionType := domain.ActionDeploy

	if client.ActiveConfigID != nil && *client.ActiveConfigID == configID {
		s.log.Info("config already active → restart",
			"client_id", clientID,
		)

		actionType = domain.ActionRestart

		if err := client.Transition(domain.ClientStatusRestarting); err != nil {
			return err
		}
	} else {
		client.ActivateConfig(configID)

		if err := client.Transition(domain.ClientStatusDeploying); err != nil {
			return err
		}

		s.log.Info("config activated",
			"client_id", clientID,
			"config_id", configID,
		)
	}

	if err := s.clientRepo.Update(ctx, client); err != nil {
		return err
	}

	action, err := s.actionService.Create(ctx, clientID, userID, actionType)
	if err != nil {
		return err
	}

	s.log.Info("deploy action created",
		"client_id", clientID,
		"action_id", action.ID,
		"type", actionType,
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

	if headers == nil {
		headers = map[string]string{}
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
		s.log.Error("failed to get client",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	config, err := s.configRepo.GetByID(ctx, configID)
	if err != nil {
		s.log.Error("failed to get config",
			"config_id", configID,
			"error", err,
		)
		return err
	}

	if config.ClientID != client.ID {
		return domain.ErrForbidden
	}

	if (client.ActiveConfigID != nil && *client.ActiveConfigID == configID) {
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
		s.log.Error("failed to get client",
			"client_id", clientID,
			"error", err,
		)
		return nil, err
	}

	oldConfig, err := s.configRepo.GetByID(ctx, configID)
	if err != nil {
		s.log.Error("failed to get config",
			"config_id", configID,
			"error", err,
		)
		return nil, err
	}

	if oldConfig.ClientID != client.ID {
		return nil, domain.ErrForbidden
	}

	configs, err := s.configRepo.ListByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	for _, c := range configs {
		if c.Version == req.Version {
			return nil, domain.ErrConfigVersionExists
		}
	}

	headers := req.Headers
	if headers == nil {
		headers = map[string]string{}
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
		s.log.Error("failed to create config",
			"config_id", newConfig.ID,
			"error", err,
		)
		return nil, err
	}

	if client.ActiveConfigID != nil && *client.ActiveConfigID == configID {

		client.ActivateConfig(newConfig.ID)

		if err := client.Transition(domain.ClientStatusDeploying); err != nil {
			return nil, err
		}

		if err := s.clientRepo.Update(ctx, client); err != nil {
			return nil, err
		}

		action, err := s.actionService.Create(ctx, clientID, userID, domain.ActionDeploy)
		if err != nil {
			return nil, err
		}

		s.log.Info("deploy triggered after config update",
			"client_id", clientID,
			"new_config_id", newConfig.ID,
			"action_id", action.ID,
		)
	}

	s.log.Info("config updated (new version created)",
		"client_id", clientID,
		"old_config_id", configID,
		"new_config_id", newConfig.ID,
	)

	return newConfig, nil
}
