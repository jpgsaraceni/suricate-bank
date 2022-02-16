package schemas

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type GetBalanceResponse struct {
	AccountID string `json:"account_id"`
	Balance   string `json:"balance"`
}

func BalanceToResponse(accountID account.ID, balance int) GetBalanceResponse {
	return GetBalanceResponse{
		AccountID: accountID.String(),
		Balance:   money.MustParseBRL(balance),
	}
}
