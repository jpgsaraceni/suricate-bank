package loginroute

import (
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/services/auth"
	"github.com/jpgsaraceni/suricate-bank/config"
)

type handler struct {
	service auth.Service
	Config  config.Config
}

type Handler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func NewHandler(cfg config.Config, s auth.Service) Handler {
	return &handler{
		service: s,
		Config:  cfg,
	}
}
