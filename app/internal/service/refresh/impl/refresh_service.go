package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/service/refresh"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type refreshService struct {
	rdb      *redis.Client
	ref_time int
}

func NewRefreshService(rdb *redis.Client, ttl int) refresh.RefreshService {
	return &refreshService{
		rdb: rdb,
		ref_time: ttl,
	}
}

func (s *refreshService) Create(ctx context.Context, userID string) (string, error) {
	token := uuid.NewString()

	err := s.rdb.Set(ctx, "refresh:"+token, userID, time.Duration(s.ref_time)*time.Second).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *refreshService) Validate(ctx context.Context, token string) (string, error) {
    userID, err := s.rdb.Get(ctx, "refresh:"+token).Result()

    if err == redis.Nil {
        return "", domain.ErrInvalidRefreshToken
    }

    if err != nil {
        return "", err
    }

    return userID, nil
}

func (s *refreshService) Delete(ctx context.Context, token string) error {
	return s.rdb.Del(ctx, "refresh:"+token).Err()
}
