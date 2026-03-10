package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"sync"
)

type InMemoryUserRepository struct {
	mu      sync.RWMutex
	storage map[string]*domain.User // key = email
}

func NewInMemoryUserRepository() repository.UserRepository {
	return &InMemoryUserRepository{
		storage: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.storage[user.Email]; exists {
		return domain.ErrUserAlreadyExists
	}

	r.storage[user.Email] = user
	return nil
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	user, ok := r.storage[email]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.storage {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}
