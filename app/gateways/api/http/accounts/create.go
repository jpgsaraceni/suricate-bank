package accountsroute

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
)

type CreateRequest struct {
	Name   string `json:"name" validate:"required"`
	Cpf    string `json:"cpf" validate:"required"`
	Secret string `json:"secret" validate:"required"`
}

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	var createRequest CreateRequest
	var response responses.Response

	defer func(r *responses.Response) {
		responses.SendJSON(w, *r)
	}(&response)

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		response = responses.BadRequest(responses.ErrInvalidRequestPayload)

		return
	}

	if createRequest.Name == "" || createRequest.Cpf == "" || createRequest.Secret == "" {
		response = responses.BadRequest(responses.ErrMissingFields)

		return
	}

	cpf := createRequest.Cpf

	if len(cpf) != 11 && len(cpf) != 14 {
		response = responses.BadRequest(responses.ErrLengthCpf)

		return
	}

	name := createRequest.Name

	if len(name) < 3 {
		response = responses.BadRequest(responses.ErrShortName)

		return
	}

	secret := createRequest.Secret

	if len(secret) < 6 {
		response = responses.BadRequest(responses.ErrShortSecret)

		return
	}

	_, err := h.usecase.Create(r.Context(), name, cpf, secret)

	if errors.Is(err, account.ErrInvalidCpf) {
		response = responses.BadRequest(responses.ErrInvalidCpf)

		return
	}

	if errors.Is(err, accountspg.ErrCpfAlreadyExists) {
		response = responses.BadRequest(responses.ErrCpfAlreadyExists)

		return
	}

	if err != nil {
		response = responses.InternalServerError(err)

		return
	}

	response = responses.Created(responses.AccountCreated)
}
