package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	errConfigureDb = errors.New("failed to configure db connection")
	errConnectDb   = errors.New("failed to connect to db")
)

// ConnectPool builds a config using the url passed as argument,
// then creates a new pool and connects using that config.
func ConnectPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)

	if err != nil {

		return nil, fmt.Errorf("%w: %s", errConfigureDb, err.Error())
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)

	if err != nil {
		return nil, fmt.Errorf("%w: %s", errConnectDb, err.Error())
	}

	migration, err := migrate.New(
		"file://../app/gateways/db/postgres/migrations",
		url)

	if err != nil {
		return nil, fmt.Errorf("could not read migration files: %s", err)
	}

	err = migration.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {

		return nil, err
	}

	return pool, nil
}
