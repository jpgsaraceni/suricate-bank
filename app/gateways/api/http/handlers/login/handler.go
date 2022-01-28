package loginroute

import (
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/services/auth"
)

type handler struct {
	service auth.Service
}

type Handler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func NewHandler(s auth.Service) Handler {
	return &handler{service: s}
}
