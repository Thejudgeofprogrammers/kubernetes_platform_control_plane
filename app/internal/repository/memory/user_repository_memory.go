package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"sync"
)

type InMemoryUserRepository struct {
	mu      sync.RWMutex
	usersByID map[string]*domain.User
	usersByEmail map[string]*domain.User
}

func NewInMemoryUserRepository() repository.UserRepository {
	return &InMemoryUserRepository{
		usersByID: make(map[string]*domain.User),
		usersByEmail: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.usersByEmail[user.Email]; exists {
		return domain.ErrUserAlreadyExists
	}

    r.usersByID[user.ID] = user
    r.usersByEmail[user.Email] = user

	return nil
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	user, ok := r.usersByEmail[email]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

    user, ok := r.usersByID[id]
    if !ok {
        return nil, domain.ErrUserNotFound
    }

    return user, nil
}
