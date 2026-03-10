package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"log"
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
	log.Println("code:", code)
	return nil
}

func (r *InMemoryEmailCodeRepository) Get(ctx context.Context, email string) (*domain.EmailCode, error) {
	r.mu.RLock()
	code, ok := r.storage[email]
	r.mu.RUnlock()

	if !ok {
		return nil, domain.ErrCodeNotFound
	}

	if time.Now().After(code.ExpiresAt) {
		r.mu.Lock()
		delete(r.storage, email)
		r.mu.Unlock()
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
