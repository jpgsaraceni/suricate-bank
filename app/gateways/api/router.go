package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	transferuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/transfer"
	accountsroute "github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/handlers/accounts"
	loginroute "github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/handlers/login"
	transfersroute "github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/handlers/transfers"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/middlewares"
	"github.com/jpgsaraceni/suricate-bank/app/services/auth"
)

func NewRouter(
	accountUC accountuc.Usecase,
	transferUC transferuc.Usecase,
	authService auth.Service,
) {
	accountsHandler := accountsroute.NewHandler(accountUC)
	transfersHandler := transfersroute.NewHandler(transferUC)
	loginHandler := loginroute.NewHandler(authService)

	r := chi.NewRouter()

	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/accounts", accountsHandler.Create)
	r.Get("/accounts", accountsHandler.Fetch)
	r.Get("/accounts/{id}/balance", accountsHandler.GetBalance)

	r.Post("/transfers", middlewares.Authorize(transfersHandler.Create))
	r.Get("/transfers", transfersHandler.Fetch)

	r.Post("/login", loginHandler.Login)

	http.ListenAndServe(":8080", r) // TODO: get port from env var
}
