package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/repository"
	health "control_plane/internal/service/health"
	"time"
)

type healthService struct {
	repo repository.ClientHealthRepostiory
	log  logger.Logger
}

func NewHealthService(repo repository.ClientHealthRepostiory, log logger.Logger) health.HealthService {
	return &healthService{
		repo: repo,
		log:  log,
	}
}

func (s *healthService) Update(ctx context.Context, clientID string, status domain.HealthStatus, message string) error {
	health := &domain.APIClientHealth{
		ClientID:  clientID,
		Status:    status,
		LastCheck: time.Now(),
		Message:   message,
	}

	s.log.Info("health update requested",
		"client_id", clientID,
		"status", status,
	)

	err := s.repo.Update(ctx, health)
	if err != nil {
		s.log.Error("health update failed",
			"client_id", clientID,
			"error", err,
		)
		return err
	}

	s.log.Info("health updated",
		"client_id", clientID,
	)

	return nil
}

func (s *healthService) Get(ctx context.Context, clientID string) (*domain.APIClientHealth, error) {
	s.log.Info("health get requested",
		"client_id", clientID,
	)

	health, err := s.repo.GetByClientID(ctx, clientID)
	if err != nil {
		s.log.Error("health get failed",
			"client_id", clientID,
			"error", err,
		)
		return nil, err
	}

	if health == nil {
		s.log.Info("health not found",
			"client_id", clientID,
		)
		return nil, nil
	}

	s.log.Info("health fetched",
		"client_id", clientID,
		"status", health.Status,
	)

	return health, nil
}

func (s *healthService) Set(clientID string, status domain.HealthStatus) {
	s.repo.Set(clientID, domain.APIClientHealth{
		ClientID:  clientID,
		Status:    status,
		LastCheck: time.Now(),
	})
}
