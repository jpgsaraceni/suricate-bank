package idempotency

import "github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"

func (s service) SetKeyValue(key string, request schema.CachedResponse) error {
	return s.repository.SetKeyValue(key, request)
}
