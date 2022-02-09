package idempotency

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

type MockRepository struct {
	OnGetCachedResponse func(ctx context.Context, key string) (schema.CachedResponse, error)
	OnCacheResponse     func(ctx context.Context, request schema.CachedResponse) error
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) GetCachedResponse(ctx context.Context, key string) (schema.CachedResponse, error) {
	return mock.OnGetCachedResponse(ctx, key)
}

func (mock MockRepository) CacheResponse(ctx context.Context, res schema.CachedResponse) error {
	return mock.OnCacheResponse(ctx, res)
}
