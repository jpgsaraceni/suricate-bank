package transferuc

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type MockDebiter struct {
	OnDebit func(id account.AccountId, amount money.Money) error
}

var _ Debiter = (*MockDebiter)(nil)

func (mock MockDebiter) Debit(id account.AccountId, amount money.Money) error {
	return mock.OnDebit(id, amount)
}

type MockCrediter struct {
	OnCredit func(id account.AccountId, amount money.Money) error
}

var _ Crediter = (*MockCrediter)(nil)

func (mock MockCrediter) Credit(id account.AccountId, amount money.Money) error {
	return mock.OnCredit(id, amount)
}
