package schemas

import (
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

type CreateRequest struct {
	Name   string `json:"name"`
	Cpf    string `json:"cpf"`
	Secret string `json:"secret"`
}

type CreateResponse struct {
	AccountId string    `json:"account_id"`
	Name      string    `json:"name"`
	Cpf       string    `json:"cpf"`
	Balance   string    `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func CreatedAccountToResponse(createdAccount account.Account) CreateResponse {
	return CreateResponse{
		AccountId: createdAccount.Id.String(),
		Name:      createdAccount.Name,
		Cpf:       createdAccount.Cpf.Masked(),
		Balance:   createdAccount.Balance.BRL(),
		CreatedAt: createdAccount.CreatedAt,
	}
}
