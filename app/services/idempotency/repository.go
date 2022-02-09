package idempotency

import "github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"

type Repository interface {
	GetCachedResponse(key string) (schema.CachedResponse, error)
	CacheResponse(key string, request schema.CachedResponse) error
}
