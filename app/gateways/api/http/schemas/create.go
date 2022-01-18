package schemas

import (
	"encoding/json"
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
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
}

func CreatedAccountToResponse(createdAccount account.Account) CreateResponse {
	return CreateResponse{
		AccountId: createdAccount.Id,
		Name:      createdAccount.Name,
		Cpf:       createdAccount.Cpf.Masked(),
		Balance:   createdAccount.Balance.BRL(),
		CreatedAt: createdAccount.CreatedAt.Local(),
	}
}

func (r CreateResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
