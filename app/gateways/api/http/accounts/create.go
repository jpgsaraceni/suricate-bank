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

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	var createRequest CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		ErrorResponse(w, err)

		return
	}

	cpf := createRequest.Cpf

	if len(cpf) != 11 && len(cpf) != 14 {
		ErrorResponse(w, ErrLengthCpf)

		return
	}

	name := createRequest.Name

	if len(name) < 3 {
		ErrorResponse(w, ErrShortName)

		return
	}

	secret := createRequest.Secret

	if len(secret) < 6 {
		ErrorResponse(w, ErrShortSecret)

		return
	}

	if _, err := h.usecase.Create(r.Context(), name, cpf, secret); err != nil {
		ErrorResponse(w, err)

		return
	}

	responses.SendJSON(w, responses.Created("successfully created"))
}
