package postgrestest

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
)

const defaultTimeout = 30

func GetTestPool() (*pgxpool.Pool, func()) {
	var dbPool *pgxpool.Pool

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=suricate",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://postgres:secret@%s/suricate?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(60) // Tell docker to hard kill the container in 60 seconds

	dockerPool.MaxWait = timeout() * time.Second
	// connects to db in container, with exponential backoff-retry,
	// because the application in the container might not be ready to accept connections yet
	if err = dockerPool.Retry(func() error {
		dbPool, err = connectTestPool(context.Background(), databaseUrl)

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// tearDown should be called to destroy container at the end of the test
	tearDown := func() {
		dbPool.Close()
		dockerPool.Purge(resource)
	}

	return dbPool, tearDown
}

func connectTestPool(ctx context.Context, databaseUrl string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseUrl)

	if err != nil {

		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)

	if err != nil {
		return nil, err
	}

	err = migrateTestDb(databaseUrl)

	if err != nil {
		return nil, err
	}

	return pool, nil
}

func migrateTestDb(databaseUrl string) error {
	migration, err := migrate.New(
		"file://../migrations",
		databaseUrl)

	if err != nil {
		return fmt.Errorf("could not read migration files: %s", err)
	}

	err = migration.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {

		return err
	}

	return nil
}

func timeout() time.Duration {
	timeoutString := os.Getenv("DOCKERTEST_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutString)

	if err != nil {

		return time.Duration(defaultTimeout)
	}

	return time.Duration(timeout)
}
