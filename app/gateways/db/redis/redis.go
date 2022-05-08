package redis

import (
	"fmt"
	"os"
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
			log.Info().Msgf("attempting to connect to redis on %s...", addr)

			c, err := redis.DialURL(os.Getenv("REDIS_URL"))
			if err != nil {
				return redis.Dial("tcp", addr)
			}

			return c, err
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
