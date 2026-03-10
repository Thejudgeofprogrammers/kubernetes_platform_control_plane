package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	health "control_plane/internal/service/health"
	"time"
)

type healthService struct {
	repo repository.ClientHealthRepostiory
}

func NewHealthService(repo repository.ClientHealthRepostiory) health.HealthService {
	return &healthService{repo: repo}
}

func (s *healthService) Update(ctx context.Context, clientID string, status domain.HealthStatus, message string) error {
	health := &domain.APIClientHealth{
		ClientID: clientID,
		Status: status,
		LastCheck: time.Now(),
		Message: message,
	}
	return s.repo.Update(ctx, health)
}

func (s *healthService) Get(ctx context.Context, clientID string) (*domain.APIClientHealth, error) {
	return s.repo.GetByClientID(ctx, clientID)
}
