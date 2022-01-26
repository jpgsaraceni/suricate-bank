package transfersroute

import (
	"net/http"

	transferuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/transfer"
)

type handler struct {
	usecase transferuc.Usecase
}

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Fetch(w http.ResponseWriter, r *http.Request)
}

func NewHandler(uc transferuc.Usecase) Handler {
	return &handler{usecase: uc}
}
