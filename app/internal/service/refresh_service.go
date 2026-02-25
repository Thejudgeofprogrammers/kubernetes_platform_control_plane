package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RefreshService struct {
	rdb      *redis.Client
	ref_time int
}

func NewRefreshService(rdb *redis.Client) *RefreshService {
	return &RefreshService{rdb: rdb}
}

func (s *RefreshService) Create(ctx context.Context, userID string) (string, error) {
	token := uuid.NewString()

	err := s.rdb.Set(ctx, "refresh:"+token, userID, time.Duration(s.ref_time)*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *RefreshService) Validate(ctx context.Context, token string) (string, error) {
	return s.rdb.Get(ctx, "refresh:"+token).Result()
}

func (s *RefreshService) Delete(ctx context.Context, token string) error {
	return s.rdb.Del(ctx, "refresh:"+token).Err()
}
