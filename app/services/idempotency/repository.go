package idempotency

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

type Repository interface {
	GetCachedResponse(ctx context.Context, key string) (schema.CachedResponse, error)
	CacheResponse(ctx context.Context, request schema.CachedResponse) error
}
