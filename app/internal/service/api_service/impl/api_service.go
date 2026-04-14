package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	apiservice "control_plane/internal/service/api_service"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type APIServiceServiceImpl struct {
	repo repository.APIServiceRepository
	log  *slog.Logger
}

func NewAPIServiceService(repo repository.APIServiceRepository, log *slog.Logger) apiservice.APIServiceService {
	return &APIServiceServiceImpl{
		repo: repo,
		log:  log,
	}
}

func (s *APIServiceServiceImpl) Create(ctx context.Context, name, baseURL, protocol string) (*domain.APIService, error) {
	s.log.Info("create api service started",
		"name", name,
		"base_url", baseURL,
		"protocol", protocol,
	)

	service := &domain.APIService{
		ID:        uuid.New().String(),
		Name:      name,
		BaseURL:   baseURL,
		Protocol:  protocol,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, service); err != nil {
		s.log.Error("failed to create api service",
			"name", name,
			"error", err,
		)
		return nil, err
	}

	s.log.Info("api service created",
		"id", service.ID,
		"name", service.Name,
	)

	return service, nil
}

func (s *APIServiceServiceImpl) List(ctx context.Context) ([]*domain.APIService, error) {
	s.log.Info("list api services started")

	list, err := s.repo.List(ctx)
	if err != nil {
		s.log.Error("failed to list api services",
			"error", err,
		)
		return nil, err
	}

	s.log.Info("api services listed",
		"count", len(list),
	)

	return list, nil
}

func (s *APIServiceServiceImpl) Delete(ctx context.Context, id string) error {
	s.log.Info("delete api service started",
		"id", id,
	)

	if id == "" {
		s.log.Warn("empty api service id")

		return errors.New("id is required")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrClientNotFound) {
			s.log.Warn("api service not found",
				"id", id,
			)
			return err
		}

		s.log.Error("failed to delete api service",
			"id", id,
			"error", err,
		)
		return err
	}

	s.log.Info("api service deleted",
		"id", id,
	)

	return nil
}

func (s *APIServiceServiceImpl) GetByID(
    ctx context.Context,
    id string,
) (*domain.APIService, error) {

    s.log.Debug("get api service",
        "api_service_id", id,
    )

    apiService, err := s.repo.GetByID(ctx, id)
    if err != nil {
        s.log.Error("failed to get api service",
            "api_service_id", id,
            "error", err,
        )
        return nil, err
    }

    if apiService == nil {
        s.log.Warn("api service not found",
            "api_service_id", id,
        )
        return nil, domain.ErrAPIServiceNotFound
    }

    return apiService, nil
}
