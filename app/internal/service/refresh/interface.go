package refresh

import "context"

type RefreshService interface {
	Create(ctx context.Context, userID string) (string, error)
	Validate(ctx context.Context, token string) (string, error)
	Delete(ctx context.Context, token string) error
}
