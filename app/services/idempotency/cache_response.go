package idempotency

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

func (s service) CacheResponse(ctx context.Context, request schema.CachedResponse) error {

	response, err := s.repository.GetCachedResponse(ctx, request.Key)

	if err != nil {

		return fmt.Errorf("%w:%s", ErrRepository, err)
	}

	if !reflect.DeepEqual(response, schema.CachedResponse{}) {

		return ErrResponseExists
	}

	err = s.repository.CacheResponse(ctx, request)

	if err != nil {

		return fmt.Errorf("%w:%s", ErrRepository, err)
	}

	return nil
}
