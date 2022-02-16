package schemas

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type GetBalanceResponse struct {
	AccountID string `json:"account_id" example:"3d368560-e8e4-4108-8bb8-9f8753db09af"`
	Balance   string `json:"balance" example:"R$10,00"`
}

func BalanceToResponse(accountID account.ID, balance int) GetBalanceResponse {
	return GetBalanceResponse{
		AccountID: accountID.String(),
		Balance:   money.MustParseBRL(balance),
	}
}
