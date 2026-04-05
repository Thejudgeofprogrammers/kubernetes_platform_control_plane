package impl

import (
	"context"
	"control_plane/internal/domain"
	"control_plane/internal/service/refresh"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type refreshService struct {
	rdb      *redis.Client
	ref_time int
	log      *slog.Logger
}

func NewRefreshService(rdb *redis.Client, ttl int, log *slog.Logger) refresh.RefreshService {
	return &refreshService{
		rdb:      rdb,
		ref_time: ttl,
		log:      log,
	}
}

func (s *refreshService) Create(ctx context.Context, userID string) (string, error) {

	s.log.Debug("refresh token create started",
		"user_id", userID,
	)

	token := uuid.NewString()

	err := s.rdb.Set(ctx, "refresh:"+token, userID, time.Duration(s.ref_time)*time.Second).Err()
	if err != nil {
		s.log.Error("failed to store refresh token",
			"user_id", userID,
			"error", err,
		)
		return "", err
	}

	s.log.Debug("refresh token created",
		"user_id", userID,
		"ttl_sec", s.ref_time,
	)

	return token, nil
}

func (s *refreshService) Validate(ctx context.Context, token string) (string, error) {

	s.log.Debug("refresh token validate started")

	userID, err := s.rdb.Get(ctx, "refresh:"+token).Result()

	if err == redis.Nil {
		s.log.Warn("invalid refresh token")
		return "", domain.ErrInvalidRefreshToken
	}

	if err != nil {
		s.log.Error("failed to get refresh token",
			"error", err,
		)
		return "", err
	}

	s.log.Debug("refresh token valid",
		"user_id", userID,
	)

	return userID, nil
}

func (s *refreshService) Delete(ctx context.Context, token string) error {
	s.log.Debug("refresh token delete started")

	err := s.rdb.Del(ctx, "refresh:"+token).Err()
	if err != nil {
		s.log.Error("failed to delete refresh token",
			"error", err,
		)
		return err
	}

	s.log.Debug("refresh token deleted")

	return nil
}
