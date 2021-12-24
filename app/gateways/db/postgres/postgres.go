package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

var errConfigureDb = errors.New("failed to configure db connection")

func ConnectPool(url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)

	if err != nil {

		return nil, errConfigureDb
	}

	pool, err := pgxpool.ConnectConfig(context.TODO(), config)
	return pool, err
}
