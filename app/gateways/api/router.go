package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

func NewRouter(
	ctx context.Context,
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

	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/accounts", middlewares.Idempotency(idempotencyService, accountsHandler.Create))
	r.Get("/accounts", accountsHandler.Fetch)
	r.Get("/accounts/{id}/balance", accountsHandler.GetBalance)

	r.With(middlewares.Authorize).Post("/transfers", transfersHandler.Create)
	r.Get("/transfers", transfersHandler.Fetch)

	r.Post("/login", loginHandler.Login)

	http.ListenAndServe(cfg.HttpServer.ServerPort(), r)
}
