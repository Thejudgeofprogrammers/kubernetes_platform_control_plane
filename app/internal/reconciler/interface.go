package reconciler

import "context"

type ReconcilerService interface {
	Run(ctx context.Context)
	Start(ctx context.Context)
}
