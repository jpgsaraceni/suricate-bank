package main

import (
	"context"
	"fmt"

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
		panic(fmt.Errorf("failed to connect to db: %w", err))
	}
	defer pgPool.Close()

	accountsRepository := accountspg.NewRepository(pgPool)
	transfersRepository := transferspg.NewRepository(pgPool)

	accountsUsecase := accountuc.NewUsecase(accountsRepository)
	transfersUsecase := transferuc.NewUsecase(transfersRepository, accountsRepository)

	authService := auth.NewService(accountsRepository)

	api.NewRouter(accountsUsecase, transfersUsecase, authService)
}
