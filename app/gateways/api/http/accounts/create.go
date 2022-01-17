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
	var response responses.Response
	response.Writer = w

	var createRequest CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		response.BadRequest(responses.ErrInvalidRequestPayload).SendJSON()

		return
	}

	if createRequest.Name == "" || createRequest.Cpf == "" || createRequest.Secret == "" {
		response.BadRequest(responses.ErrMissingFields).SendJSON()

		return
	}

	cpf := createRequest.Cpf

	if len(cpf) != 11 && len(cpf) != 14 {
		response.BadRequest(responses.ErrLengthCpf).SendJSON()

		return
	}

	name := createRequest.Name

	if len(name) < 3 {
		response.BadRequest(responses.ErrShortName).SendJSON()

		return
	}

	secret := createRequest.Secret

	if len(secret) < 6 {
		response.BadRequest(responses.ErrShortSecret).SendJSON()

		return
	}

	_, err := h.usecase.Create(r.Context(), name, cpf, secret)

	if errors.Is(err, account.ErrInvalidCpf) {
		response.BadRequest(responses.ErrInvalidCpf).SendJSON()

		return
	}

	if errors.Is(err, accountspg.ErrCpfAlreadyExists) {
		response.BadRequest(responses.ErrCpfAlreadyExists).SendJSON()

		return
	}

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}
	response.Created(responses.AccountCreated).SendJSON()
}
