package schemas

import (
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

type FetchAccountsResponse struct {
	Accounts []FetchedAccount `json:"accounts"`
}

type FetchedAccount struct {
	AccountId string    `json:"account_id"`
	Name      string    `json:"name"`
	Cpf       string    `json:"cpf"`
	Balance   string    `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func AccountsToResponse(accountList []account.Account) FetchAccountsResponse {
	accountResponse := make([]FetchedAccount, 0, len(accountList))
	for _, account := range accountList {
		accountResponse = append(accountResponse, FetchedAccount{
			AccountId: account.Id.String(),
			Name:      account.Name,
			Cpf:       account.Cpf.Masked(),
			Balance:   account.Balance.BRL(),
			CreatedAt: account.CreatedAt,
		})
	}
	return FetchAccountsResponse{Accounts: accountResponse}
}
