package mock

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/orchestrator"
	"control_plane/internal/service/health"
	"log/slog"
)

type MockOrchestrator struct {
	healthService health.HealthService
	log           *slog.Logger
}

func NewMockOrchestrator(
	healthService health.HealthService,
	log *slog.Logger,
) orchestrator.Orchestrator {
	return &MockOrchestrator{
		healthService: healthService,
		log:           log,
	}
}

func (m *MockOrchestrator) Deploy(
	ctx context.Context,
	client *domain.APIClient,
	config *domain.APIClientConfig,
) error {
	m.log.Info("deploy started",
		"client_id", client.ID,
		"config_id", config.ID,
	)

	err := m.healthService.Update(
		ctx,
		client.ID,
		domain.HealthHealthy,
		"",
	)

	if err != nil {
		m.log.Error("deploy failed",
			"client_id", client.ID,
			"error", err,
		)
		return err
	}

	m.log.Info("deploy completed",
		"client_id", client.ID,
	)

	return nil
}

func (m *MockOrchestrator) Restart(
	ctx context.Context,
	clientID string,
) error {
	m.log.Info("restart started",
		"client_id", clientID,
	)

	err := m.healthService.Update(
		ctx,
		clientID,
		domain.HealthHealthy,
		"",
	)

	if err != nil {
		m.log.Error("restart failed",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	m.log.Info("restart completed",
		"client_id", clientID,
	)

	return nil
}

func (m *MockOrchestrator) Delete(
	ctx context.Context,
	clientID string,
) error {
	m.log.Info("delete started",
		"client_id", clientID,
	)

	err := m.healthService.Update(
		ctx,
		clientID,
		domain.HealthUnhealthy,
		"client deleted",
	)

	if err != nil {
		m.log.Error("delete failed",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	m.log.Info("delete completed",
		"client_id", clientID,
	)

	return nil
}
