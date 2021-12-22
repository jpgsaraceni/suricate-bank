package account

import "github.com/jpgsaraceni/suricate-bank/app/vos/money"

type MockRepository struct {
	OnCreate        func(account *Account) error
	OnGetBalance    func(id AccountId) (int, error)
	OnFetch         func() ([]Account, error)
	OnGetById       func(id AccountId) (Account, error)
	OnCreditAccount func(account *Account, amount money.Money) error
	OnDebitAccount  func(account *Account, amount money.Money) error
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) Create(account *Account) error {
	return mock.OnCreate(account)
}

func (mock MockRepository) GetBalance(id AccountId) (int, error) {
	return mock.OnGetBalance(id)
}

func (mock MockRepository) Fetch() ([]Account, error) {
	return mock.OnFetch()
}

func (mock MockRepository) GetById(id AccountId) (Account, error) {
	return mock.OnGetById(id)
}

func (mock MockRepository) CreditAccount(account *Account, amount money.Money) error {
	return mock.OnCreditAccount(account, amount)
}

func (mock MockRepository) DebitAccount(account *Account, amount money.Money) error {
	return mock.OnDebitAccount(account, amount)
}
