package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"

	"github.com/jpgsaraceni/suricate-bank/config"
)

const (
	maxIdle     = 3
	idleTimeout = 240
)

func ConnectPool(cfg config.Config) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout * time.Second,
		Dial: func() (redis.Conn, error) {
			url := cfg.GetRedisURL()
			log.Info().Msgf("attempting to connect to redis on %s...", url)

			c, err := redis.DialURL(url)
			if err != nil {
				return redis.Dial("tcp", url)
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
