package memory

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/repository"
	"log/slog"
	"sync"
	"time"
)

type InMemoryEmailCodeRepository struct {
	mu      sync.RWMutex
	storage map[string]*domain.EmailCode // key = email

	log *slog.Logger
}

func NewInMemoryEmailCodeRepository(log *slog.Logger) repository.EmailCodeRepository {
	return &InMemoryEmailCodeRepository{
		storage: make(map[string]*domain.EmailCode),
		log:     log,
	}
}

func (r *InMemoryEmailCodeRepository) Save(ctx context.Context, code *domain.EmailCode) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[code.Email] = code

	r.log.Info("email code saved",
		"email", code.Email,
		"expires_at", code.ExpiresAt,
	)

	return nil
}

func (r *InMemoryEmailCodeRepository) Get(ctx context.Context, email string) (*domain.EmailCode, error) {
	r.mu.RLock()
	code, ok := r.storage[email]
	r.mu.RUnlock()

	if !ok {
		r.log.Error("code not found",
			"email", email,
		)
		return nil, domain.ErrCodeNotFound
	}

	if time.Now().After(code.ExpiresAt) {
		r.mu.Lock()
		delete(r.storage, email)
		r.mu.Unlock()

		r.log.Info("code expired and deleted",
			"email", email,
		)

		return nil, domain.ErrCodeExpired
	}

	copyCode := *code

	r.log.Info("code fetched",
		"email", email,
	)

	return &copyCode, nil
}

func (r *InMemoryEmailCodeRepository) Delete(ctx context.Context, email string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.storage, email)

	r.log.Info("code deleted",
		"email", email,
	)

	return nil
}
