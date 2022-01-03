package account

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type MockRepository struct {
	OnCreate        func(ctx context.Context, account *Account) error
	OnGetBalance    func(ctx context.Context, id AccountId) (int, error)
	OnFetch         func(ctx context.Context) ([]Account, error)
	OnGetById       func(ctx context.Context, id AccountId) (Account, error)
	OnCreditAccount func(ctx context.Context, account *Account, amount money.Money) error
	OnDebitAccount  func(ctx context.Context, account *Account, amount money.Money) error
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) Create(ctx context.Context, account *Account) error {
	return mock.OnCreate(ctx, account)
}

func (mock MockRepository) GetBalance(ctx context.Context, id AccountId) (int, error) {
	return mock.OnGetBalance(ctx, id)
}

func (mock MockRepository) Fetch(ctx context.Context) ([]Account, error) {
	return mock.OnFetch(ctx)
}

func (mock MockRepository) GetById(ctx context.Context, id AccountId) (Account, error) {
	return mock.OnGetById(ctx, id)
}

func (mock MockRepository) CreditAccount(ctx context.Context, account *Account, amount money.Money) error {
	return mock.OnCreditAccount(ctx, account, amount)
}

func (mock MockRepository) DebitAccount(ctx context.Context, account *Account, amount money.Money) error {
	return mock.OnDebitAccount(ctx, account, amount)
}
