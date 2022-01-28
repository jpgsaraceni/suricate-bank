package schemas

import (
	"errors"
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

type CreateAccountRequest struct {
	Name   string `json:"name"`
	Cpf    string `json:"cpf"`
	Secret string `json:"secret"`
}

type CreateAccountResponse struct {
	AccountId string    `json:"account_id"`
	Name      string    `json:"name"`
	Cpf       string    `json:"cpf"`
	Balance   string    `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func CreatedAccountToResponse(createdAccount account.Account) CreateAccountResponse {
	return CreateAccountResponse{
		AccountId: createdAccount.Id.String(),
		Name:      createdAccount.Name,
		Cpf:       createdAccount.Cpf.Masked(),
		Balance:   createdAccount.Balance.BRL(),
		CreatedAt: createdAccount.CreatedAt,
	}
}

const (
	minNameLength = 3
	maxNameLength = 100

	minSecretLength = 6

	rawCpfLength    = 11
	maskedCpfLength = 14
)

func (r CreateAccountRequest) Validate(response responses.Response) (account.Account, responses.Response) {
	if r.Name == "" || r.Cpf == "" || r.Secret == "" {

		return account.Account{}, response.BadRequest(responses.ErrMissingFieldsAccountPayload)
	}

	if len(r.Name) < minNameLength || len(r.Name) > maxNameLength {

		return account.Account{}, response.BadRequest(responses.ErrLengthName)
	}

	if len(r.Secret) < minSecretLength {

		return account.Account{}, response.BadRequest(responses.ErrShortSecret)
	}

	if len(r.Cpf) != rawCpfLength && len(r.Cpf) != maskedCpfLength {

		return account.Account{}, response.BadRequest(responses.ErrLengthCpf)
	}

	accountInstance, err := account.NewAccount(r.Name, r.Cpf, r.Secret)

	if err != nil {
		if errors.Is(err, account.ErrInvalidCpf) {

			return account.Account{}, response.BadRequest(responses.ErrInvalidCpf)
		}

		return account.Account{}, response.InternalServerError(err)
	}

	return accountInstance, response
}
