package redis

import (
	"encoding/json"
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

var errKeyExists = errors.New("key already exists in redis")

func (r Repository) CacheResponse(response schema.CachedResponse) error {
	conn := r.pool.Get()
	defer conn.Close()

	j, err := json.Marshal(response)

	if err != nil {

		return err
	}

	reply, err := conn.Do("SETNX", response.Key, j)

	if reply.(int64) == 0 {

		return errKeyExists
	}

	if err != nil {

		return err
	}

	return nil
}
