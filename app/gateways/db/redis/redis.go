package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
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
			log.Info().Msg(fmt.Sprintf("attempting to connect to redis on %s...\n", addr))

			return redis.Dial("tcp", addr)
		},
	}

	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil { // check if redis server is responsive
		return nil, fmt.Errorf("failed to ping redis server: %w", err)
	}
	log.Info().Msg("successfully connected to redis server")

	return pool, nil
}
