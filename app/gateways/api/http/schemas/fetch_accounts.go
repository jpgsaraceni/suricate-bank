package schemas

import (
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

type FetchAccountsResponse struct {
	Accounts []FetchedAccount `json:"accounts"`
}

type FetchedAccount struct {
	AccountID string    `json:"account_id" example:"5738eda2-49f5-4702-83e4-b87b18cf0d31"`
	Name      string    `json:"name" example:"Zé do Caroço"`
	Cpf       string    `json:"cpf" example:"220.614.460-35"`
	Balance   string    `json:"balance" example:"R$10,00"`
	CreatedAt time.Time `json:"created_at" example:"2022-01-28T19:39:04.585238-03:00"`
}

func AccountsToResponse(accountList []account.Account) FetchAccountsResponse {
	accountResponse := make([]FetchedAccount, 0, len(accountList))
	for _, account := range accountList {
		accountResponse = append(accountResponse, FetchedAccount{
			AccountID: account.ID.String(),
			Name:      account.Name,
			Cpf:       account.Cpf.Masked(),
			Balance:   account.Balance.BRL(),
			CreatedAt: account.CreatedAt,
		})
	}

	return FetchAccountsResponse{Accounts: accountResponse}
}
