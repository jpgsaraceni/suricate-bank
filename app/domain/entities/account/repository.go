package account

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type Repository interface {
	Create(ctx context.Context, account *Account) error
	GetBalance(ctx context.Context, id AccountId) (int, error)
	Fetch(ctx context.Context) ([]Account, error)
	GetById(ctx context.Context, id AccountId) (Account, error)
	CreditAccount(ctx context.Context, account *Account, amount money.Money) error
	DebitAccount(ctx context.Context, account *Account, amount money.Money) error
}
