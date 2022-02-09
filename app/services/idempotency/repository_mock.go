package idempotency

import "github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"

type MockRepository struct {
	OnGetCachedResponse func(key string) (schema.CachedResponse, error)
	OnCacheResponse     func(request schema.CachedResponse) error
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) GetCachedResponse(key string) (schema.CachedResponse, error) {
	return mock.OnGetCachedResponse(key)
}

func (mock MockRepository) CacheResponse(res schema.CachedResponse) error {
	return mock.OnCacheResponse(res)
}
