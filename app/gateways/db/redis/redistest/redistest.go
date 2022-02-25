package redistest

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rs/zerolog/log"
)

const (
	containerTimeout = 60
	poolTimeout      = 30
	maxIdle          = 3
	idleTimeout      = 240
)

func GetTestPool() (*redis.Pool, func()) {
	var dbPool *redis.Pool

	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Panic().Err(err).Msg("Could not connect to docker")
	}

	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "6-alpine",
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Panic().Err(err).Msg("Could not start resource")
	}

	hostAndPort := resource.GetHostPort("6379/tcp")

	log.Info().Msgf("Connecting to redis on: %s", hostAndPort)

	if err = resource.Expire(containerTimeout); err != nil { // Tell docker to hard kill the container in 60 seconds
		panic(err)
	}

	dockerPool.MaxWait = poolTimeout * time.Second
	// connects to db in container, with exponential backoff-retry,
	// because the application in the container might not be ready to accept connections yet
	if err = dockerPool.Retry(func() error {
		dbPool = &redis.Pool{
			MaxIdle:     maxIdle,
			IdleTimeout: idleTimeout * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", hostAndPort)
			},
		}

		dbConn := dbPool.Get()
		_, err = dbConn.Do("PING")

		return err
	}); err != nil {
		log.Panic().Err(err).Msg("Could not connect to docker")
	}

	tearDown := func() {
		if err = dockerPool.Purge(resource); err != nil {
			panic(err)
		}
	}

	return dbPool, tearDown
}
