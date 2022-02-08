package redis

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

var errKeyExists = errors.New("key already exists in redis")

func (r Repository) SetKeyValue(key string, res responses.Response) error {
	conn := r.pool.Get()
	defer conn.Close()

	payloadJson, err := json.Marshal(res.Payload)

	if err != nil {

		return fmt.Errorf("failed to marshal response payload: %s", err)
	}

	reply, err := conn.Do("HSETNX", key, "status", res.Status)
	if err != nil {

		return fmt.Errorf("redis HSETNX command error: %s", err)
	}
	if reply == 0 {

		return errKeyExists
	}

	reply, err = conn.Do("HSETNX", key, "payload", payloadJson)
	if err != nil {

		return fmt.Errorf("redis HSETNX command error: %s", err)
	}
	if reply == 0 {

		return errKeyExists
	}

	// redis unstable version:
	// reply, err = conn.Do("HSETNX", key, "status", res.status, "payload", payloadJson)
	// if err != nil {

	// 	return fmt.Errorf("redis HSETNX command error: %s", err)
	// }
	// if reply == 0 {

	// 	return errKeyExists
	// }

	return nil
}
