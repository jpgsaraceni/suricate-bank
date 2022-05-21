package idempotency

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
	"github.com/jpgsaraceni/suricate-bank/config"
)

type Repository interface {
	GetCachedResponse(ctx context.Context, key string) (schema.CachedResponse, error)
	CacheResponse(ctx context.Context, request schema.CachedResponse) error
	Lock(ctx context.Context, cfg config.Config, key string) error
}
