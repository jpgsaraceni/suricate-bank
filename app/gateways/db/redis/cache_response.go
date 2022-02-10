package redis

import (
	"context"
	"encoding/json"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

func (r Repository) CacheResponse(ctx context.Context, response schema.CachedResponse) error {
	conn := r.pool.Get()
	defer conn.Close()

	j, err := json.Marshal(response)

	if err != nil {

		return err
	}

	reply, err := conn.Do("SET", response.Key, j, "XX")

	if reply == nil {

		return errKeyNotFound
	}

	if err != nil {

		return err
	}

	return nil
}
