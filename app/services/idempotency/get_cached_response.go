package idempotency

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

func (s service) GetCachedResponse(ctx context.Context, key string) (schema.CachedResponse, error) {
	return s.repository.GetCachedResponse(key)
}
