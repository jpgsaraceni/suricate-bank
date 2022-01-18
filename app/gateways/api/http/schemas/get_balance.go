package schemas

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type GetBalanceResponse struct {
	AccountId account.AccountId `json:"account_id"`
	Balance   string            `json:"balance"`
}

func BalanceToResponse(accountId account.AccountId, balance int) GetBalanceResponse {
	return GetBalanceResponse{
		AccountId: accountId,
		Balance:   money.MustParseBRL(balance),
	}
}
