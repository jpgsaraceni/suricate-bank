package accountsroute

import (
	"encoding/json"
	"fmt"
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

	cpf := createRequest.Cpf

	if len(cpf) != 11 && len(cpf) != 14 {

		return responses.BadRequest(fmt.Errorf("invalid cpf length"), ErrLengthCpf)
	}

	name := createRequest.Name

	if len(name) < 3 {

		return responses.BadRequest(fmt.Errorf("name too short"), ErrShortName)
	}

	secret := createRequest.Secret

	if len(secret) < 6 {

		return responses.BadRequest(fmt.Errorf("secret too short"), ErrShortSecret)
	}

	if _, err := h.Usecase.Create(r.Context(), name, cpf, secret); err != nil {
		return responses.InternalServerError(err, ErrInternalServerError)
	}

	return responses.Response{}
}