package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/repository"
	"sort"
	"sync"
)

type InMemoryClientActionRepository struct {
	storage []*domain.APIClientAction
	mu      sync.RWMutex

	log logger.Logger
}

func NewInMemoryClientActionRepository(log logger.Logger) repository.ClientActionRepository {
	return &InMemoryClientActionRepository{
		storage: []*domain.APIClientAction{},
		log:     log,
	}
}

func (r *InMemoryClientActionRepository) Create(ctx context.Context, action *domain.APIClientAction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage = append(r.storage, action)

	r.log.Info("action created",
		"id", action.ID,
		"client_id", action.ClientID,
		"status", action.Status,
	)

	return nil
}

func (r *InMemoryClientActionRepository) ListByClientID(ctx context.Context, clientID string) ([]*domain.APIClientAction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.APIClientAction, 0, len(r.storage))

	for _, action := range r.storage {
		if action.ClientID == clientID {
			copyAction := *action
			result = append(result, &copyAction)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	r.log.Info("list actions by client",
		"client_id", clientID,
		"count", len(result),
	)

	return result, nil
}

func (r *InMemoryClientActionRepository) GetPending(ctx context.Context) ([]*domain.APIClientAction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.APIClientAction

	for _, action := range r.storage {
		if action.Status == domain.ActionPending {
			copyAction := *action
			result = append(result, &copyAction)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	r.log.Info("get pending actions",
		"count", len(result),
	)

	return result, nil
}

func (r *InMemoryClientActionRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status domain.ActionStatus,
	errMsg *string,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, action := range r.storage {
		if action.ID == id {
			action.Status = status
			action.Error = errMsg

			r.log.Info("action status updated",
				"id", id,
				"status", status,
				"error", errMsg,
			)

			return nil
		}
	}

	r.log.Error("action not found",
		"id", id,
	)

	return domain.ErrClientNotFound
}
