package reconciler

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/orchestrator"
	"control_plane/internal/reconciler"
	"control_plane/internal/repository"
	health "control_plane/internal/service/health"
	metics "control_plane/internal/service/metric"
	"fmt"
	"log"

	"time"
)

type reconcilerService struct {
	actionRepo     repository.ClientActionRepository
	clientRepo     repository.ClientRepository
	apiServiceRepo repository.APIServiceRepository
	orchestrator   orchestrator.Orchestrator
	configRepo     repository.ClientConfigRepository
	log            logger.Logger
	healthSrv      health.HealthService
	metricsRepo    repository.MetricsRepository
	metricsService metics.MetricsService
	baseURL        string
}

func NewReconciler(
	actionRepo repository.ClientActionRepository,
	clientRepo repository.ClientRepository,
	apiServiceRepo repository.APIServiceRepository,
	orchestrator orchestrator.Orchestrator,
	configRepo repository.ClientConfigRepository,
	log logger.Logger,
	healthS health.HealthService,
	metricsRepo repository.MetricsRepository,
	metricsService metics.MetricsService,
	baseURL string,
) reconciler.ReconcilerService {
	return &reconcilerService{
		actionRepo:     actionRepo,
		clientRepo:     clientRepo,
		apiServiceRepo: apiServiceRepo,
		orchestrator:   orchestrator,
		configRepo:     configRepo,
		log:            log,
		healthSrv:      healthS,
		metricsRepo:    metricsRepo,
		metricsService: metricsService,
		baseURL:        baseURL,
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
			r.syncState(ctx)
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
	if action.Status != domain.ActionPending {
		return
	}

	r.log.Info("action started",
		"action_id", action.ID,
		"client_id", action.ClientID,
		"type", action.Type,
	)

	_ = r.actionRepo.UpdateStatus(ctx, action.ID, domain.ActionRunning, nil)

	client, err := r.clientRepo.GetByID(ctx, action.ClientID)
	if err != nil {
		r.failAction(ctx, action, err)
		return
	}

	var execErr error

	switch action.Type {

	case domain.ActionStop:
		r.log.Info("action stop", "client_id", client.ID)

		execErr = r.orchestrator.Delete(ctx, client.ID)

	case domain.ActionRestart:
		r.log.Info("action restart", "client_id", client.ID)

		if err := client.Transition(domain.ClientStatusRestarting); err != nil {
			execErr = err
			break
		}

		execErr = r.orchestrator.Restart(ctx, client.ID)
		if execErr != nil {
			break
		}

		execErr = client.Transition(domain.ClientStatusRunning)

	case domain.ActionDelete:
		r.log.Info("action delete", "client_id", client.ID)

		execErr = r.orchestrator.Delete(ctx, client.ID)
		if execErr != nil {
			break
		}

		execErr = client.Transition(domain.ClientStatusDisabled)

	case domain.ActionDeploy:
		r.log.Info("action deploy", "client_id", client.ID)

		if client.ActiveConfigID == nil {
			execErr = fmt.Errorf("no active config")
			break
		}

		config, err := r.getConfig(ctx, *client.ActiveConfigID)
		if err != nil {
			execErr = err
			break
		}

		execErr = r.orchestrator.Deploy(ctx, client, config)
		if execErr != nil {
			break
		}

		execErr = client.Transition(domain.ClientStatusRunning)
	}

	if execErr != nil {
		r.failAction(ctx, action, execErr)
		return
	}

	if err := r.clientRepo.Update(ctx, client); err != nil {
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

func (r *reconcilerService) syncState(ctx context.Context) {
	clients, _, err := r.clientRepo.List(ctx, "", 1000, 0)
	if err != nil {
		r.log.Error("failed to list clients", "error", err)
		return
	}

	for _, c := range clients {

		if c.GetStatus() == domain.ClientStatusDisabled {
			continue
		}

		metrics, err := r.metricsService.Collect(ctx, r.baseURL, c.ID)
		if err == nil {
			for _, m := range metrics {
				_ = r.metricsRepo.Save(ctx, m)
			}
		}

		r.orchestrator.CheckHealth(ctx, c.ID)

		health, err := r.healthSrv.Get(ctx, c.ID)
		if err != nil || health == nil {
			continue
		}

		oldStatus := c.GetStatus()

		switch health.Status {
		case domain.HealthHealthy:
			if c.GetStatus() == domain.ClientStatusStopping {
				_ = c.Transition(domain.ClientStatusStopped)
			}

		case domain.HealthDegraded:
			r.log.Warn("client degraded", "client_id", c.ID)
		// 	if err := c.Transition(domain.ClientStatusDeploying); err != nil {
		// 		r.log.Warn("failed to transition to deploing",
		// 			"client_id", c.ID,
		// 			"current_status", c.GetStatus(),
		// 			"error", err,
		// 		)
		// 	}

		case domain.HealthUnhealthy:
			if c.GetStatus() == domain.ClientStatusStopping {
				if err := c.Transition(domain.ClientStatusStopped); err != nil {
					r.log.Warn("failed to sync stopped",
						"client_id", c.ID,
						"status", c.GetStatus(),
						"error", err,
					)
				}
			}
		}

		if oldStatus != c.GetStatus() {
			_ = r.clientRepo.Update(ctx, c)
		}
	}
}

func (r *reconcilerService) failAction(ctx context.Context, action *domain.APIClientAction, err error) {
	msg := err.Error()

	r.log.Error("action failed",
		"action_id", action.ID,
		"client_id", action.ClientID,
		"error", err,
	)

	_ = r.actionRepo.UpdateStatus(ctx, action.ID, domain.ActionFailed, &msg)
}

func (r *reconcilerService) getConfig(ctx context.Context, id string) (*domain.APIClientConfig, error) {
	r.log.Debug("get config",
		"config_id", id,
	)

	return r.configRepo.GetByID(ctx, id)
}
