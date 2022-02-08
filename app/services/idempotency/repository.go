package idempotency

import "github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"

type Repository interface {
	GetKeyValue(key string) (schema.CachedResponse, error)
	SetKeyValue(key string, request schema.CachedResponse) error
}
