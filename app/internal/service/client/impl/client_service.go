package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/orchestrator"
	"control_plane/internal/repository"
	"control_plane/internal/service/action"
	"control_plane/internal/service/client"

	"github.com/google/uuid"
)

type clientService struct {
	repo       repository.ClientRepository
	orch       orchestrator.Orchestrator
	configRepo repository.ClientConfigRepository
	actionSrv  action.ActionService
	log        logger.Logger
}

func NewClientService(
	repo repository.ClientRepository,
	orch orchestrator.Orchestrator,
	configRepo repository.ClientConfigRepository,
	actionSrv action.ActionService,
	log logger.Logger,
) client.ClientService {
	return &clientService{
		repo:       repo,
		orch:       orch,
		configRepo: configRepo,
		actionSrv:  actionSrv,
		log:        log,
	}
}

func (s *clientService) List(ctx context.Context, status string, limit, offset int) ([]domain.APIClient, int, error) {
	s.log.Debug("list clients",
		"status", status,
		"limit", limit,
		"offset", offset,
	)

	clients, total, err := s.repo.List(ctx, status, limit, offset)
	if err != nil {
		s.log.Error("failed to list clients",
			"status", status,
			"error", err,
		)
		return nil, 0, err
	}

	result := make([]domain.APIClient, 0, len(clients))

	for _, c := range clients {
		result = append(result, *c)
	}

	s.log.Debug("clients listed",
		"count", len(result),
		"total", total,
	)

	return result, total, nil
}

func (s *clientService) Create(ctx context.Context, userID, name, apiServiceID, description string) (*domain.APIClient, error) {
	s.log.Info("client create started",
		"user_id", userID,
		"name", name,
		"api_service_id", apiServiceID,
	)

	id := uuid.NewString()
	client := domain.NewAPIClient(id, name, description, apiServiceID)

	if err := s.repo.Create(ctx, client); err != nil {
		s.log.Error("failed to create client",
			"client_id", id,
			"error", err,
		)
		return nil, err
	}

	s.log.Info("client created",
		"client_id", client.ID,
		"user_id", userID,
	)

	return client, nil
}

func (s *clientService) GetByID(ctx context.Context, id string) (*domain.APIClient, error) {

	s.log.Debug("get client",
		"client_id", id,
	)

	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get client",
			"client_id", id,
			"error", err,
		)
		return nil, err
	}

	return client, nil
}

func (s *clientService) Restart(ctx context.Context, userID, id, reason string) error {
	s.log.Info("client restart started",
		"client_id", id,
		"user_id", userID,
		"reason", reason,
	)

	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get client",
			"client_id", id,
			"error", err,
		)
		return err
	}

	s.log.Info("client found",
		"client_id", client.ID,
		"status", client.GetStatus(),
	)

	if !client.CanRestart() {
		return domain.ErrInvalidStateTransition
	}

	if err := client.Transition(domain.ClientStatusRestarting); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, client); err != nil {
		return err
	}

	action, err := s.actionSrv.Create(ctx, id, userID, domain.ActionRestart)
	if err != nil {
		s.log.Error("failed to create restart action",
			"client_id", id,
			"user_id", userID,
			"error", err,
		)
		return err
	}

	s.log.Info("client restart scheduled",
		"client_id", id,
		"action_id", action.ID,
	)

	return nil
}

func (s *clientService) Delete(ctx context.Context, userID, id string) error {
	s.log.Info("client delete started",
		"client_id", id,
		"user_id", userID,
	)

	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get client",
			"client_id", id,
			"error", err,
		)
		return err
	}

	if err := client.Transition(domain.ClientStatusDeleting); err != nil {
		s.log.Error("invalid state transition",
			"client_id", id,
			"error", err,
		)
		return err
	}

	if err := s.repo.Update(ctx, client); err != nil {
		s.log.Error("failed to update client",
			"client_id", id,
			"error", err,
		)
		return err
	}

	action, err := s.actionSrv.Create(ctx, id, userID, domain.ActionDelete)
	if err != nil {
		return err
	}

	s.log.Info("client delete scheduled",
		"client_id", id,
		"action_id", action.ID,
	)

	return nil
}

func (s *clientService) Start(ctx context.Context, userID, clientID string) error {
	s.log.Info("client start requested",
		"client_id", clientID,
		"user_id", userID,
	)

	client, err := s.repo.GetByID(ctx, clientID)
	if err != nil {
		return err
	}

	if !client.CanStart() {
		return domain.ErrInvalidStateTransition
	}

	if err := client.Transition(domain.ClientStatusDeploying); err != nil {
		return err
	}	

	if err := s.repo.Update(ctx, client); err != nil {
		return err
	}

	action, err := s.actionSrv.Create(ctx, clientID, userID, domain.ActionDeploy)
	if err != nil {
		return err
	}

	s.log.Info("client start scheduled",
		"client_id", clientID,
		"action_id", action.ID,
	)

	return nil
}

func (s *clientService) Stop(
	ctx context.Context,
	userID string,
	clientID string,
) error {
	s.log.Info("stop client started",
		"client_id", clientID,
		"user_id", userID,
	)

	client, err := s.repo.GetByID(ctx, clientID)
	if err != nil {
		return err
	}

	if !client.CanStop() {
		return domain.ErrInvalidStateTransition
	}

	if err := client.Transition(domain.ClientStatusStopping); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, client); err != nil {
		return err
	}

	action, err := s.actionSrv.Create(ctx, clientID, userID, domain.ActionStop)
	if err != nil {
		return err
	}

	s.log.Info("stop action created",
		"client_id", clientID,
		"action_id", action.ID,
	)

	return nil
}
