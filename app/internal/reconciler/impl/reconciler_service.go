package reconciler

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/orchestrator"
	"control_plane/internal/reconciler"
	"control_plane/internal/repository"
	"log"
	"log/slog"
	"time"
)

type reconcilerService struct {
	actionRepo   repository.ClientActionRepository
	clientRepo   repository.ClientRepository
	orchestrator orchestrator.Orchestrator
	configRepo   repository.ClientConfigRepository
	log          *slog.Logger
}

func NewReconciler(
	actionRepo repository.ClientActionRepository,
	clientRepo repository.ClientRepository,
	orchestrator orchestrator.Orchestrator,
	configRepo repository.ClientConfigRepository,
	log *slog.Logger,
) reconciler.ReconcilerService {
	return &reconcilerService{
		actionRepo:   actionRepo,
		clientRepo:   clientRepo,
		orchestrator: orchestrator,
		configRepo:   configRepo,
		log:          log,
	}
}

func (r *reconcilerService) Run(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	r.log.Info("reconciler started")

	for {
		select {
		case <-ctx.Done():
			log.Println("reconciler stopped")
			return

		case <-ticker.C:
			r.process(ctx)
		}
	}
}

func (r *reconcilerService) process(ctx context.Context) {
	r.log.Debug("reconciler tick")

	actions, err := r.actionRepo.GetPending(ctx)
	if err != nil {
		r.log.Error("failed to get pending actions",
			"error", err,
		)
		return
	}

	if len(actions) == 0 {
		r.log.Debug("no pending actions")
		return
	}

	r.log.Info("processing actions",
		"count", len(actions),
	)

	for _, action := range actions {
		r.handleAction(ctx, action)
	}
}

func (r *reconcilerService) handleAction(ctx context.Context, action *domain.APIClientAction) {

	r.log.Info("action started",
		"action_id", action.ID,
		"client_id", action.ClientID,
		"type", action.Type,
	)

	_ = r.actionRepo.UpdateStatus(ctx, action.ID, domain.ActionRunning, nil)

	client, err := r.clientRepo.GetByID(ctx, action.ClientID)
	if err != nil {
		msg := err.Error()

		r.log.Error("failed to get client",
			"action_id", action.ID,
			"client_id", action.ClientID,
			"error", err,
		)

		_ = r.actionRepo.UpdateStatus(ctx, action.ID, domain.ActionFailed, &msg)
		return
	}

	switch action.Type {

	case domain.ActionRestart:
		r.log.Info("action restart",
			"client_id", client.ID,
		)

		err = r.orchestrator.Restart(ctx, client.ID)

		if err == nil {
			_ = client.Transition(domain.ClientStatusRunning)
		}

	case domain.ActionDelete:
		r.log.Info("action delete",
			"client_id", client.ID,
		)

		err = r.orchestrator.Delete(ctx, client.ID)

		if err == nil {
			_ = client.Transition(domain.ClientStatusDisabled)
		}

	case domain.ActionDeploy:
		r.log.Info("action deploy",
			"client_id", client.ID,
		)

		if client.ActiveConfigID == nil {
			msg := "no active config"

			r.log.Error("deploy failed: no active config",
				"action_id", action.ID,
				"client_id", client.ID,
			)

			_ = r.actionRepo.UpdateStatus(ctx, action.ID, domain.ActionFailed, &msg)
			return
		}

		config, err2 := r.getConfig(ctx, *client.ActiveConfigID)
		if err2 != nil {
			msg := err2.Error()

			r.log.Error("failed to get config",
				"action_id", action.ID,
				"client_id", client.ID,
				"config_id", *client.ActiveConfigID,
				"error", err2,
			)

			_ = r.actionRepo.UpdateStatus(ctx, action.ID, domain.ActionFailed, &msg)
			return
		}

		err = r.orchestrator.Deploy(ctx, client, config)

		if err == nil {
			_ = client.Transition(domain.ClientStatusRunning)
		}
	}

	if err != nil {
		msg := err.Error()

		r.log.Error("action failed",
			"action_id", action.ID,
			"client_id", client.ID,
			"type", action.Type,
			"error", err,
		)

		_ = r.actionRepo.UpdateStatus(ctx, action.ID, domain.ActionFailed, &msg)
		return
	}

	err = r.clientRepo.Update(ctx, client)
	if err != nil {
		r.log.Error("failed to update client",
			"client_id", client.ID,
			"error", err,
		)
	}

	_ = r.actionRepo.UpdateStatus(ctx, action.ID, domain.ActionSuccess, nil)

	r.log.Info("action completed",
		"action_id", action.ID,
		"client_id", client.ID,
		"type", action.Type,
	)
}

func (r *reconcilerService) getConfig(ctx context.Context, id string) (*domain.APIClientConfig, error) {
	r.log.Debug("get config",
		"config_id", id,
	)

	return r.configRepo.GetByID(ctx, id)
}
