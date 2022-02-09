package idempotency

import "github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"

func (s service) CacheResponse(key string, request schema.CachedResponse) error {
	return s.repository.CacheResponse(key, request)
}
