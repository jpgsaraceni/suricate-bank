package postgrestest

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
)

var dbPool *pgxpool.Pool
var testContext = context.Background()

// TestMain creates runs a docker container of PostgreSQL to run
// integration tests.
func GetTestPool() *pgxpool.Pool {
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

	resource.Expire(60) // Tell docker to hard kill the container in 60 seconds

	dockerPool.MaxWait = 10 * time.Second
	// connects to db in container, with exponential backoff-retry,
	// because the application in the container might not be ready to accept connections yet
	if err = dockerPool.Retry(func() error {
		dbPool, err = postgres.ConnectPool(testContext, databaseUrl)

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// migration, err := os.ReadFile("../migrations/000001_accounts.up.sql")
	migration, err := migrate.New(
		"file://../migrations",
		databaseUrl)

	if err != nil {
		log.Fatalf("Could not read migration files: %s", err)
	}

	if err := migration.Up(); err != nil {
		log.Fatal(err)
	}

	// creates accounts table in db
	// if _, err := dbPool.Exec(context.Background(), string(migration)); err != nil {
	// 	log.Fatalf("Could not run migration: %s", err)
	// }

	//Run tests
	// code := m.Run()

	// dbPool.Close()

	return dbPool

	// if err := dockerPool.Purge(resource); err != nil {
	// 	log.Fatalf("Could not purge resource: %s", err)
	// }

	// os.Exit(code)
}

// truncateAccounts clears the accounts table so tests are independent
func Truncate(pool *pgxpool.Pool, table string) error {
	_, err := pool.Exec(testContext, fmt.Sprintf("TRUNCATE %s", table))

	if err != nil {

		return err
	}

	return nil
}
