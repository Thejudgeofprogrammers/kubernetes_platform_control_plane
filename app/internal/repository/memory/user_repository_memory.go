package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/logger"
	"control_plane/internal/repository"
	"sync"
)

type InMemoryUserRepository struct {
	mu           sync.RWMutex
	usersByID    map[string]*domain.User
	usersByEmail map[string]*domain.User

	log logger.Logger
}

func NewInMemoryUserRepository(log logger.Logger) repository.UserRepository {
	return &InMemoryUserRepository{
		usersByID:    make(map[string]*domain.User),
		usersByEmail: make(map[string]*domain.User),
		log:          log,
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.usersByEmail[user.Email]; exists {
		r.log.Error("user already exists", "email", user.Email)
		return domain.ErrUserAlreadyExists
	}

	r.usersByID[user.ID] = user
	r.usersByEmail[user.Email] = user

	r.log.Info("user created",
		"id", user.ID,
		"email", user.Email,
	)

	return nil
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.usersByEmail[email]
	if !ok {
		r.log.Error("user not found by email", "email", email)
		return nil, domain.ErrUserNotFound
	}

	copyUser := *user

	r.log.Info("user fetched by email",
		"id", user.ID,
		"email", email,
	)

	return &copyUser, nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.usersByID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	copyUser := *user

	r.log.Info("user fetched by id",
		"id", id,
	)

	return &copyUser, nil
}

func (r *InMemoryUserRepository) List(ctx context.Context) ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.User, 0, len(r.usersByID))

	for _, u := range r.usersByID {
		result = append(result, *u)
	}

	r.log.Info("list users",
		"count", len(result),
	)

	return result, nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.usersByID[id]
	if !ok {
		r.log.Error("delete failed: user not found", "id", id)
		return domain.ErrUserNotFound
	}

	delete(r.usersByID, id)
	delete(r.usersByEmail, user.Email)

	r.log.Info("user deleted",
		"id", id,
		"email", user.Email,
	)

	return nil
}

func (r *InMemoryUserRepository) UpdateRole(ctx context.Context, id string, role string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.usersByID[id]
	if !ok {
		r.log.Error("update role failed: user not found", "id", id)
		return domain.ErrUserNotFound
	}

	user.Role = domain.AccessRole(role)

	r.log.Info("user role updated",
		"id", id,
		"role", role,
	)

	return nil
}
