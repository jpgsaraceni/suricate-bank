package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

func (r Repository) GetKeyValue(key string) (responses.Response, error) {
	conn := r.pool.Get()
	defer conn.Close()

	reply, err := redis.Values(conn.Do("HMGET", key, "status", "payload"))

	if err != nil {

		return responses.Response{}, fmt.Errorf("redis HMGET command error: %s", err)
	}

	response, err := responses.Unmarshal(reply)
	// var response responses.Response
	// _, err = redis.Scan(reply, &response)

	if err != nil {

		return responses.Response{}, fmt.Errorf("failed to unmarshal reply from redis: %s", err)
	}

	return response, nil
}
