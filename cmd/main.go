package main

import (
	"context"
	"log"

	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	transferuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	transferspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/transfers"
	"github.com/jpgsaraceni/suricate-bank/app/services/auth"
	"github.com/jpgsaraceni/suricate-bank/config"
)

func main() {
	ctx := context.TODO()

	cfg := config.ReadConfig(".env")

	pgPool, err := postgres.ConnectPool(ctx, cfg.Postgres.Url())
	if err != nil {
		panic(err)
	}

	log.Printf("\033[34m---- HAPPY BANKING ----\033[37m\n")

	defer pgPool.Close()

	// connect to redis
	// defer redis close

	accountsRepository := accountspg.NewRepository(pgPool)
	transfersRepository := transferspg.NewRepository(pgPool)
	// instantiate redis repository

	accountsUsecase := accountuc.NewUsecase(accountsRepository)
	transfersUsecase := transferuc.NewUsecase(transfersRepository, accountsRepository)

	authService := auth.NewService(accountsRepository)
	// instantiate idempotency service

	api.NewRouter(ctx, *cfg, accountsUsecase, transfersUsecase, authService)
}
