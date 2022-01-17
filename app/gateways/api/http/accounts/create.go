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

func (h handler) Create(w http.ResponseWriter, r *http.Request) error {
	var createRequest CreateRequest
	var response responses.Response

	response.Writer = w

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {

		return response.BadRequest(responses.ErrInvalidRequestPayload).SendJSON()
	}

	if createRequest.Name == "" || createRequest.Cpf == "" || createRequest.Secret == "" {

		return response.BadRequest(responses.ErrMissingFields).SendJSON()
	}

	cpf := createRequest.Cpf

	if len(cpf) != 11 && len(cpf) != 14 {

		return response.BadRequest(responses.ErrLengthCpf).SendJSON()
	}

	name := createRequest.Name

	if len(name) < 3 {

		return response.BadRequest(responses.ErrShortName).SendJSON()
	}

	secret := createRequest.Secret

	if len(secret) < 6 {

		return response.BadRequest(responses.ErrShortSecret).SendJSON()
	}

	_, err := h.usecase.Create(r.Context(), name, cpf, secret)

	if errors.Is(err, account.ErrInvalidCpf) {

		return response.BadRequest(responses.ErrInvalidCpf).SendJSON()
	}

	if errors.Is(err, accountspg.ErrCpfAlreadyExists) {

		return response.BadRequest(responses.ErrCpfAlreadyExists).SendJSON()
	}

	if err != nil {

		return response.InternalServerError(err).SendJSON()
	}

	return response.Created(responses.AccountCreated).SendJSON()
}
