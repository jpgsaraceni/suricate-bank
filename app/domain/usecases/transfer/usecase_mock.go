package transferuc

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type MockDebiter struct {
	OnDebit func(ctx context.Context, id account.AccountId, amount money.Money) error
}

var _ Debiter = (*MockDebiter)(nil)

func (mock MockDebiter) Debit(ctx context.Context, id account.AccountId, amount money.Money) error {
	return mock.OnDebit(ctx, id, amount)
}

type MockCrediter struct {
	OnCredit func(ctx context.Context, id account.AccountId, amount money.Money) error
}

var _ Crediter = (*MockCrediter)(nil)

func (mock MockCrediter) Credit(ctx context.Context, id account.AccountId, amount money.Money) error {
	return mock.OnCredit(ctx, id, amount)
}
