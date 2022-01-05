package accountsroute

import (
	"encoding/json"
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

type CreateRequest struct {
	Name   string `json:"name" validate:"required"`
	Cpf    string `json:"cpf" validate:"required"`
	Secret string `json:"secret" validate:"required"`
}

func (h Handler) Create(r *http.Request) responses.Response {
	var createRequest CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		return responses.BadRequest(err, ErrMissingFields)
	}

	return responses.Response{}
}
