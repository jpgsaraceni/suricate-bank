package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func ConnectPool(addr string) (*redis.Conn, error) {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {

			return redis.Dial("tcp", addr)
		},
	}

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING") // check if redis server is responsive

	if err != nil {

		return nil, err
	}

	return &conn, nil
}
