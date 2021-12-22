package account

import "github.com/jpgsaraceni/suricate-bank/app/vos/money"

type Repository interface {
	Create(account *Account) error
	GetBalance(id AccountId) (int, error)
	Fetch() ([]Account, error)
	GetById(id AccountId) (Account, error)
	CreditAccount(account *Account, amount money.Money) error
	DebitAccount(account *Account, amount money.Money) error
}
