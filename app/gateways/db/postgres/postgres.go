package postgres

import (
	"context"
	"errors"
	"fmt"

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

	return pool, nil
}
