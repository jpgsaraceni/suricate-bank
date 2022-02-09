package idempotency

import "github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"

type MockRepository struct {
	OnGetCachedResponse func(key string) (schema.CachedResponse, error)
	OnCacheResponse     func(key string, request schema.CachedResponse) error
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) GetCachedResponse(key string) (schema.CachedResponse, error) {
	return mock.OnGetCachedResponse(key)
}

func (mock MockRepository) CacheResponse(key string, res schema.CachedResponse) error {
	return mock.OnCacheResponse(key, res)
}
