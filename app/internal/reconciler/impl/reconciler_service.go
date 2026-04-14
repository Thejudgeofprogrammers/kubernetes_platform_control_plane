package reconciler

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/orchestrator"
	"control_plane/internal/reconciler"
	"control_plane/internal/repository"
	health "control_plane/internal/service/health"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type reconcilerService struct {
	actionRepo     repository.ClientActionRepository
	clientRepo     repository.ClientRepository
	apiServiceRepo repository.APIServiceRepository
	orchestrator   orchestrator.Orchestrator
	configRepo     repository.ClientConfigRepository
	log            *slog.Logger
	healthSrv      health.HealthService
}

func NewReconciler(
	actionRepo repository.ClientActionRepository,
	clientRepo repository.ClientRepository,
	apiServiceRepo repository.APIServiceRepository,
	orchestrator orchestrator.Orchestrator,
	configRepo repository.ClientConfigRepository,
	log *slog.Logger,
	healthS health.HealthService,
) reconciler.ReconcilerService {
	return &reconcilerService{
		actionRepo:     actionRepo,
		clientRepo:     clientRepo,
		apiServiceRepo: apiServiceRepo,
		orchestrator:   orchestrator,
		configRepo:     configRepo,
		log:            log,
		healthSrv:      healthS,
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

		if err := client.Transition(domain.ClientStatusRestarting); err != nil {
			r.log.Warn("failed to transition to restarting",
				"client_id", client.ID,
				"error", err,
			)
		}

		err = r.orchestrator.Restart(ctx, client.ID)

		if err == nil {
			if err2 := client.Transition(domain.ClientStatusRunning); err2 != nil {
				r.log.Warn("failed to transition to running",
					"client_id", client.ID,
					"error", err2,
				)
			}
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

func (r *reconcilerService) Start(ctx context.Context) {
	go func() {
		for {
			clients, _, err := r.clientRepo.List(ctx, "", 1000, 0) // Тут доделать нужно
			if err != nil {
				r.log.Error("failed to list clients", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, c := range clients {
				r.orchestrator.CheckHealth(ctx, c.ID)

				health, err := r.healthSrv.Get(ctx, c.ID)
				if err != nil {
					r.log.Error("failed to get health", "client_id", c.ID, "error", err)
					continue
				}

				if health == nil {
					continue
				}

				oldStatus := c.GetStatus()

				switch health.Status {
				case domain.HealthHealthy:
					if c.GetStatus() == domain.ClientStatusCreated {
						_ = c.Transition(domain.ClientStatusDeploying)
					}

					_ = c.Transition(domain.ClientStatusRunning)

				case domain.HealthDegraded:
					_ = c.Transition(domain.ClientStatusDeploying)

				case domain.HealthUnhealthy:
					_ = c.Transition(domain.ClientStatusStopped)
				}

				if oldStatus != c.GetStatus() {
					if err := r.clientRepo.Update(ctx, c); err != nil {
						r.log.Error("failed to update client status",
							"client_id", c.ID,
							"error", err,
						)
					}
				}
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

func ParseMetrics(raw string, clientID string) []domain.Metric {
	lines := strings.Split(raw, "\n")
	var result []domain.Metric

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			continue
		}

		name := parts[0]
		value, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}

		result = append(result, domain.Metric{
			ClientID:  clientID,
			Name:      name,
			Value:     value,
			CreatedAt: time.Now(),
		})
	}

	return result
}
