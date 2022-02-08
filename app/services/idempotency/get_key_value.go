package idempotency

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

func (s service) GetKeyValue(ctx context.Context, key string) (schema.CachedResponse, error) {
	return s.repository.GetKeyValue(key)
}
