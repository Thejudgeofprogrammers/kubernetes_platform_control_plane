package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"sync"
	"time"
)

type InMemoryEmailCodeRepository struct {
	mu      sync.RWMutex
	storage map[string]*domain.EmailCode // key = email
}

func NewInMemoryEmailCodeRepository() repository.EmailCodeRepository {
	return &InMemoryEmailCodeRepository{
		storage: make(map[string]*domain.EmailCode),
	}
}

func (r *InMemoryEmailCodeRepository) Save(ctx context.Context, code *domain.EmailCode) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[code.Email] = code
	return nil
}

func (r *InMemoryEmailCodeRepository) Get(ctx context.Context, email string) (*domain.EmailCode, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	code, ok := r.storage[email]
	if !ok {
		return nil, domain.ErrCodeNotFound
	}

	if time.Now().After(code.ExpiresAt) {
		return nil, domain.ErrCodeExpired
	}

	return code, nil
}

func (r *InMemoryEmailCodeRepository) Delete(ctx context.Context, email string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.storage, email)
	return nil
}
