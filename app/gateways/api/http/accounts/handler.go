package accountsroute

import (
	"net/http"

	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
)

// Handler will be used to bind all handlers for the /accounts route and access usecase.
type handler struct {
	usecase accountuc.Usecase
}

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request) error
	GetBalance(w http.ResponseWriter, r *http.Request) error
	Fetch(w http.ResponseWriter, r *http.Request) error
}

func NewHandler(uc accountuc.Usecase) Handler {
	return &handler{usecase: uc}
}
