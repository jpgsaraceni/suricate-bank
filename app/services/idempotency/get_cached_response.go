package idempotency

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

func (s service) GetCachedResponse(ctx context.Context, key string) (schema.CachedResponse, error) {
	response, err := s.repository.GetCachedResponse(ctx, key)
	if err != nil {
		return response, fmt.Errorf("%w:%s", ErrRepository, err)
	}

	return response, nil
}
