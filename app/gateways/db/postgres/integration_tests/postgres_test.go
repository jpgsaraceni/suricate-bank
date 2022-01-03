package postgres_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
)

var (
	dbPool      *pgxpool.Pool
	testContext = context.Background()
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=suricate",
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
	databaseUrl := fmt.Sprintf("postgres://suricate:secret@%s/suricate?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	dockerPool.MaxWait = 10 * time.Second
	if err = dockerPool.Retry(func() error {
		dbPool, err = postgres.ConnectPool(databaseUrl)

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	migration, err := os.ReadFile("../migrations/000001_accounts.up.sql")

	if err != nil {
		log.Fatalf("Could not read migration file: %s", err)
	}

	if _, err := dbPool.Exec(context.Background(), string(migration)); err != nil {
		log.Fatalf("Could not run migration: %s", err)
	}

	//Run tests
	code := m.Run()

	dbPool.Close()

	if err := dockerPool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func truncateAccounts() error {
	_, err := dbPool.Exec(testContext, "TRUNCATE accounts")

	if err != nil {

		return err
	}

	return nil
}
