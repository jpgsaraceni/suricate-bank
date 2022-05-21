package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddlewares "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
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
	loginHandler := loginroute.NewHandler(cfg, authService)

	r := chi.NewRouter()

	r.Use(chiMiddlewares.Timeout(requestTimeout * time.Second))
	r.Use(chiMiddlewares.Recoverer)
	r.Use(middlewares.RequestID) // wrapper for chi requestID middleware
	r.Use(middlewares.RequestLogger)

	r.With(
		middlewares.Idempotency(cfg, idempotencyService),
	).Post("/accounts", accountsHandler.Create)
	r.Get("/accounts", accountsHandler.Fetch)
	r.Get("/accounts/{id}/balance", accountsHandler.GetBalance)

	r.With(
		middlewares.Authorize(cfg),
		middlewares.Idempotency(cfg, idempotencyService),
	).Post("/transfers", transfersHandler.Create)
	r.Get("/transfers", transfersHandler.Fetch)

	r.Post("/login", loginHandler.Login)

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "swagger/index.html", http.StatusMovedPermanently)
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	if err := http.ListenAndServe(cfg.GetHTTPPort(), r); err != nil {
		log.Panic().Stack().Msgf("failed to listen and serve: %s", err)
	}
}
