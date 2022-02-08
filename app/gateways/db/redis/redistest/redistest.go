package redistest

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func GetTestPool() (*redis.Pool, func()) {
	var dbPool *redis.Pool

	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
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
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("6379/tcp")

	log.Println("Connecting to redis on: ", hostAndPort)

	resource.Expire(60) // Tell docker to hard kill the container in 60 seconds

	dockerPool.MaxWait = 30 * time.Second
	// connects to db in container, with exponential backoff-retry,
	// because the application in the container might not be ready to accept connections yet
	if err = dockerPool.Retry(func() error {
		dbPool = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {

				return redis.Dial("tcp", hostAndPort)
			},
		}

		dbConn := dbPool.Get()
		_, err := dbConn.Do("PING")

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	tearDown := func() {
		dockerPool.Purge(resource)
	}

	return dbPool, tearDown
}
