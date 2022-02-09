package redis

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

func ConnectPool(addr string) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			log.Printf("attempting to connect to redis on %s...\n", addr)

			return redis.Dial("tcp", addr)
		},
	}

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING") // check if redis server is responsive

	if err != nil {

		return nil, err
	}
	log.Println("successfully connected to redis server")

	return pool, nil
}
