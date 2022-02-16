package redis

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	maxIdle     = 3
	idleTimeout = 240
)

func ConnectPool(addr string) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			log.Printf("attempting to connect to redis on %s...\n", addr)

			return redis.Dial("tcp", addr)
		},
	}

	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil { // check if redis server is responsive
		return nil, err
	}
	log.Println("successfully connected to redis server")

	return pool, nil
}
