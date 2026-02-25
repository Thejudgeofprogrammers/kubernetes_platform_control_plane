package mock

import (
	"context"
	"control_plane/internal/domain"
	"log"
)

type MockOrchestrator struct {}

func NewMockOrchestrator() *MockOrchestrator {
	return &MockOrchestrator{}
}

func (m *MockOrchestrator) Deploy(
	ctx context.Context,
	client *domain.APIClient,
	config *domain.APIClientConfig,
) error {
	log.Printf("Deploy client %s with config %s", client.ID, config.ID)
	return nil
}

func (m *MockOrchestrator) Restart(
	ctx context.Context,
	clientID string,
) error {
	log.Printf("Restarting client %s", clientID)
	return nil
}

func (m *MockOrchestrator) Delete(
	ctx context.Context,
	clientID string,
) error {
	log.Printf("Deleting client %s", clientID)
	return nil
}