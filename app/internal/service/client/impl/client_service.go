package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/orchestrator"
	"control_plane/internal/repository"
	"control_plane/internal/service/audit"
	"control_plane/internal/service/client"

	"github.com/google/uuid"
)

type clientService struct {
	repo  repository.ClientRepository
	orch  orchestrator.Orchestrator
	audit audit.AuditService
}

func NewClientService(repo repository.ClientRepository, orch orchestrator.Orchestrator, audit audit.AuditService) client.ClientService {
	return &clientService{
		repo:  repo,
		audit: audit,
		orch:  orch,
	}
}

func (s *clientService) List(ctx context.Context, status string, limit, offset int) ([]domain.APIClient, int, error) {
	clients, total, err := s.repo.List(ctx, status, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]domain.APIClient, 0, len(clients))

	for _, c := range clients {
		result = append(result, *c)
	}

	return result, total, nil
}

func (s *clientService) Create(ctx context.Context, userID, name, apiServiceID, description string) (*domain.APIClient, error) {
	id := uuid.NewString()

	client := domain.NewAPIClient(id, name, description, apiServiceID)

	err := s.repo.Create(ctx, client)
	if err != nil {
		return nil, err
	}

	err = s.audit.Log(ctx, client.ID, userID, domain.ActionCreate)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s *clientService) GetByID(ctx context.Context, id string) (*domain.APIClient, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *clientService) Restart(ctx context.Context, userID, id, reason string) error {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := client.Transition(domain.ClientStatusRestarting); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, client); err != nil {
		return err
	}

	if err := s.orch.Restart(ctx, client.ID); err != nil {
		return err
	}

	return s.audit.Log(ctx, id, userID, domain.ActionRestart)
}

func (s *clientService) Delete(ctx context.Context, userID, id string) error {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := client.Transition(domain.ClientStatusDeleting); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, client); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return s.audit.Log(ctx, id, userID, domain.ActionDelete)
}
