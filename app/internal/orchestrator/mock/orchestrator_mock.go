package mock

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/orchestrator"
	"control_plane/internal/service/health"
	"log"
)

type MockOrchestrator struct {
	healthService health.HealthService
}

func NewMockOrchestrator(
	healthService health.HealthService,
) orchestrator.Orchestrator {
	return &MockOrchestrator{
		healthService: healthService,
	}
}

func (m *MockOrchestrator) Deploy(
	ctx context.Context,
	client *domain.APIClient,
	config *domain.APIClientConfig,
) error {
	log.Printf("Deploy client %s with config %s", client.ID, config.ID)
	
	err := m.healthService.Update(
		ctx,
		client.ID,
		domain.HealthHealthy,
		"",
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *MockOrchestrator) Restart(
	ctx context.Context,
	clientID string,
) error {
	log.Printf("Restarting client %s", clientID)

	err := m.healthService.Update(
		ctx,
		clientID,
		domain.HealthHealthy,
		"",
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *MockOrchestrator) Delete(
	ctx context.Context,
	clientID string,
) error {
	log.Printf("Deleting client %s", clientID)

	err := m.healthService.Update(
		ctx,
		clientID,
		domain.HealthUnhealthy,
		"client deleted",
	)

	if err != nil {
		return err
	}

	return nil
}