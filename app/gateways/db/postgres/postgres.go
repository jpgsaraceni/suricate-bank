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

func ConnectPool(url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)

	if err != nil {

		return nil, fmt.Errorf("%w: %s", errConfigureDb, err.Error())
	}

	pool, err := pgxpool.ConnectConfig(context.TODO(), config)

	if err != nil {
		return nil, fmt.Errorf("%w: %s", errConnectDb, err.Error())
	}

	return pool, nil
}
