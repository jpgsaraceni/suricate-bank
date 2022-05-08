package main

import (
	"context"

	"github.com/rs/zerolog/log"

	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	transferuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	transferspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/transfers"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/redis"
	"github.com/jpgsaraceni/suricate-bank/app/infrastructure/logging"
	"github.com/jpgsaraceni/suricate-bank/app/services/auth"
	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency"
	"github.com/jpgsaraceni/suricate-bank/config"
)

// @title Suricate Bank API
// @version 0.2.0
// @description Suricate Bank is an api that creates accounts and transfers money between them.
// @description ### Authorization
// @description To create a transfer (`POST /transfer`) you will need to add an Authorization header
// @description to your request, in the format **Authorization: Bearer YOUR_TOKEN**. You can set this
// @description clicking on the authorize button and entering "Bearer YOUR_TOKEN". You can get your
// @description token from the login response.
// @description ### Idempotent Requests
// @description Create transfer and account routes (`POST /transfer` and `POST /account`) support
// @description idempotent requests (you will always get the same response for the same
// @description request, without creating duplicates). To use, just set an Idempotency-Key on your
// @description request (any string, for example a UUID).

// @contact.name João Saraceni
// @contact.url https://www.linkedin.com/in/joaosaraceni/
// @contact.email jpgome@id.uff.br

// @license.name MIT
// @license.url https://github.com/jpgsaraceni/suricate-bank/blob/main/LICENSE

// @securityDefinitions.apikey Access token
// @in header
// @name Authorization
func main() {
	ctx := context.Background()

	cfg := config.ReadConfig(".env")

	logging.InitZerolog(cfg.Log.Level)

	pgPool, err := postgres.ConnectPool(ctx, cfg.Postgres.URL())
	if err != nil {
		log.Panic().Stack().Err(err).Msg("")
	}

	defer pgPool.Close()

	redisPool, err := redis.ConnectPool(cfg.Redis.URL())
	if err != nil {
		log.Warn().Msgf("failed to connect to idempotency server:%s", err)
	}

	defer redisPool.Close()

	log.Info().Msg("\033[34m---- HAPPY BANKING ----\033[37m")

	accountsRepository := accountspg.NewRepository(pgPool)
	transfersRepository := transferspg.NewRepository(pgPool)
	idempotencyRepository := redis.NewRepository(redisPool)

	accountsUsecase := accountuc.NewUsecase(accountsRepository)
	transfersUsecase := transferuc.NewUsecase(transfersRepository, accountsRepository)

	authService := auth.NewService(accountsRepository)
	idemppotencyService := idempotency.NewService(idempotencyRepository)

	// docs.SwaggerInfo.Host = cfg.HTTPServer.HostAndPort()

	api.NewRouter(ctx, *cfg, accountsUsecase, transfersUsecase, authService, idemppotencyService)
}
