package redis

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

var errType = errors.New("failed to convert redis reply to []byte")

func (r Repository) GetCachedResponse(ctx context.Context, key string) (schema.CachedResponse, error) {
	conn := r.pool.Get()
	defer conn.Close()

	var response schema.CachedResponse

	reply, err := conn.Do("GET", key)

	if err != nil {

		return response, err
	}

	if reply == nil { // key does not exist

		return response, nil
	}

	replyBytes, ok := reply.([]byte)

	if !ok {

		return response, errType
	}

	err = json.Unmarshal(replyBytes, &response)

	if err != nil {

		return response, err
	}

	return response, nil
}
