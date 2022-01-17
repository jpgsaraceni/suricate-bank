package accountsroute

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
)

type CreateRequest struct {
	Name   string `json:"name"`
	Cpf    string `json:"cpf"`
	Secret string `json:"secret"`
}

type CreateResponse struct {
	AccountId account.AccountId `json:"account_id"`
	Name      string            `json:"name"`
	Cpf       string            `json:"cpf"`
	Balance   string            `json:"balance"`
	CreatedAt time.Time         `josn:"created_at"`
	Success   bool              `json:"success"`
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

	createdAccount, err := h.usecase.Create(r.Context(), name, cpf, secret)

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

	response.Payload = CreateResponse{
		AccountId: createdAccount.Id,
		Name:      createdAccount.Name,
		Cpf:       createdAccount.Cpf.Masked(),
		Balance:   createdAccount.Balance.BRL(),
		CreatedAt: createdAccount.CreatedAt,
		Success:   true,
	}

	response.Created(responses.AccountCreated).SendJSON()
}
