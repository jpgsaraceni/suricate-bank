package accountsroute

import accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"

// Handler will be used to bind all handlers for the /accounts route,
type Handler struct {
	Usecase accountuc.Usecase
}
