package schemas

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type GetBalanceResponse struct {
	AccountId string `json:"account_id"`
	Balance   string `json:"balance"`
}

func BalanceToResponse(accountId account.AccountId, balance int) GetBalanceResponse {
	return GetBalanceResponse{
		AccountId: accountId.String(),
		Balance:   money.MustParseBRL(balance),
	}
}
