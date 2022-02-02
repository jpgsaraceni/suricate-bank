package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	errConfigureDb = errors.New("failed to configure db connection")
	errConnectDb   = errors.New("failed to connect to db")
)

// ConnectPool builds a config using the url passed as argument,
// then creates a new pool and connects using that config.
func ConnectPool(ctx context.Context, databaseUrl string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseUrl)

	if err != nil {

		return nil, fmt.Errorf("%w: %s", errConfigureDb, err.Error())
	}

	log.Printf("attempting to connect to postgres on %s...\n", databaseUrl)
	pool, err := pgxpool.ConnectConfig(ctx, config)

	if err != nil {
		return nil, fmt.Errorf("%w: %s", errConnectDb, err.Error())
	}

	log.Printf("successfully connected \n running migrations...\n")
	err = Migrate(databaseUrl)

	if err != nil {
		return nil, err
	}

	return pool, nil
}
