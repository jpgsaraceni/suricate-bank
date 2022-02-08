package redis

import (
	"encoding/json"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

func (r Repository) SetKeyValue(key string, res responses.Response) error {
	conn := r.pool.Get()
	defer conn.Close()

	payloadJson, err := json.Marshal(res.Payload)

	if err != nil {

		return fmt.Errorf("failed to marshal response payload: %s", err)
	}

	_, err = conn.Do("HSET", key, "status", res.Status, "payload", payloadJson)
	if err != nil {

		return fmt.Errorf("redis HSET command error: %s", err)
	}

	return nil
}
