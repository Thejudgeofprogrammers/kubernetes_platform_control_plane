package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/orchestrator"
	"control_plane/internal/repository"
	"control_plane/internal/service/action"
	"control_plane/internal/service/audit"
	"control_plane/internal/service/client"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type clientService struct {
	repo       repository.ClientRepository
	orch       orchestrator.Orchestrator
	audit      audit.AuditService
	configRepo repository.ClientConfigRepository
	actionSrv  action.ActionService
	actionRepo repository.ClientActionRepository
	log        *slog.Logger
}

func NewClientService(
	repo repository.ClientRepository,
	orch orchestrator.Orchestrator,
	audit audit.AuditService,
	configRepo repository.ClientConfigRepository,
	actionSrv action.ActionService,
	actionRepo repository.ClientActionRepository,
	log *slog.Logger,
) client.ClientService {
	return &clientService{
		repo:       repo,
		audit:      audit,
		orch:       orch,
		configRepo: configRepo,
		actionSrv:  actionSrv,
		actionRepo: actionRepo,
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

	if err := s.audit.Log(ctx, client.ID, userID, domain.ActionCreate); err != nil {
		s.log.Error("failed to write audit log",
			"client_id", client.ID,
			"user_id", userID,
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

	if err := s.repo.Delete(ctx, id); err != nil {
		s.log.Error("failed to delete client",
			"client_id", id,
			"error", err,
		)
		return err
	}

	if err := s.audit.Log(ctx, id, userID, domain.ActionDelete); err != nil {
		s.log.Error("failed to write audit log",
			"client_id", id,
			"user_id", userID,
			"error", err,
		)
		return err
	}

	s.log.Info("client deleted",
		"client_id", id,
	)

	return nil
}

func (s *clientService) Start(ctx context.Context, clientID string) error {
	s.log.Info("client start started",
		"client_id", clientID,
	)

	client, err := s.repo.GetByID(ctx, clientID)
	if err != nil {
		s.log.Error("failed to get client",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	if !client.CanStart() {
		s.log.Warn("client cannot start",
			"client_id", clientID,
			"status", client.GetStatus(),
		)
		return domain.ErrInvalidStateTransition
	}

	if err := client.Transition(domain.ClientStatusDeploying); err != nil {
		s.log.Error("failed to transition to deploying",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	if err := s.repo.Update(ctx, client); err != nil {
		s.log.Error("failed to update client",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	if client.ActiveConfigID == nil {
		s.log.Error("no active config",
			"client_id", clientID,
		)
		return domain.ErrConfigNotFound
	}

	config, err := s.configRepo.GetByID(ctx, *client.ActiveConfigID)
	if err != nil {
		s.log.Error("failed to get config",
			"client_id", clientID,
			"config_id", *client.ActiveConfigID,
			"error", err,
		)
		return err
	}

	s.log.Info("orchestrator deploy",
		"client_id", clientID,
		"config_id", config.ID,
	)

	if err := s.orch.Deploy(ctx, client, config); err != nil {
		s.log.Error("deploy failed",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	if err := client.Transition(domain.ClientStatusRunning); err != nil {
		s.log.Error("failed to transition to running",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	userID, _ := ctx.Value("user_id").(string)

	if err := s.audit.Log(ctx, clientID, userID, domain.ActionStart); err != nil {
		s.log.Error("failed to write audit log",
			"client_id", clientID,
			"user_id", userID,
			"error", err,
		)
	}

	if err := s.repo.Update(ctx, client); err != nil {
		s.log.Error("failed to persist running state",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	s.log.Info("client started",
		"client_id", clientID,
	)

	return nil
}

func (s *clientService) Stop(
	ctx context.Context,
	clientID string,
) error {

	userID, _ := ctx.Value("userID").(string)

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

	action := &domain.APIClientAction{
		ID:        uuid.NewString(),
		ClientID:  clientID,
		UserID:    userID,
		Type:      domain.ActionStop,
		CreatedAt: time.Now(),
	}

	if err := s.actionRepo.Create(ctx, action); err != nil {
		return err
	}

	s.log.Info("stop action created",
		"client_id", clientID,
		"action_id", action.ID,
	)

	return nil
}
