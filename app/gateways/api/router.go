package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	transferuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/transfer"
	accountsroute "github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/handlers/accounts"
	loginroute "github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/handlers/login"
	transfersroute "github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/handlers/transfers"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/middlewares"
	"github.com/jpgsaraceni/suricate-bank/app/services/auth"
	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency"
	"github.com/jpgsaraceni/suricate-bank/config"
)

const requestTimeout = 60

func NewRouter(
	_ context.Context,
	cfg config.Config,
	accountUC accountuc.Usecase,
	transferUC transferuc.Usecase,
	authService auth.Service,
	idempotencyService idempotency.Service,
) {
	accountsHandler := accountsroute.NewHandler(accountUC)
	transfersHandler := transfersroute.NewHandler(transferUC)
	loginHandler := loginroute.NewHandler(authService)

	r := chi.NewRouter()

	r.Use(middleware.Timeout(requestTimeout * time.Second))

	r.With(
		middlewares.Idempotency(idempotencyService),
	).Post("/accounts", accountsHandler.Create)
	r.Get("/accounts", accountsHandler.Fetch)
	r.Get("/accounts/{id}/balance", accountsHandler.GetBalance)

	r.With(
		middlewares.Authorize,
		middlewares.Idempotency(idempotencyService),
	).Post("/transfers", transfersHandler.Create)
	r.Get("/transfers", transfersHandler.Fetch)

	r.Post("/login", loginHandler.Login)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	if err := http.ListenAndServe(cfg.HTTPServer.ServerPort(), r); err != nil {
		log.Fatalf("failed to listen and serve: %s", err)
	}
}
