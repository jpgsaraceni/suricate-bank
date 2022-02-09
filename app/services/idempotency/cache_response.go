package idempotency

import (
	"fmt"
	"reflect"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

func (s service) CacheResponse(key string, request schema.CachedResponse) error {

	response, err := s.repository.GetCachedResponse(key)

	if err != nil {

		return fmt.Errorf("%w:%s", ErrRepository, err)
	}

	if !reflect.DeepEqual(response, schema.CachedResponse{}) {

		return ErrResponseExists
	}

	err = s.repository.CacheResponse(key, request)

	if err != nil {

		return fmt.Errorf("%w:%s", ErrRepository, err)
	}

	return nil
}
