package account

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type MockRepository struct {
	OnCreate        func(ctx context.Context, account Account) (Account, error)
	OnGetBalance    func(ctx context.Context, id ID) (int, error)
	OnFetch         func(ctx context.Context) ([]Account, error)
	OnGetByID       func(ctx context.Context, id ID) (Account, error)
	OnGetByCpf      func(ctx context.Context, cpf cpf.Cpf) (Account, error)
	OnCreditAccount func(ctx context.Context, id ID, amount money.Money) error
	OnDebitAccount  func(ctx context.Context, id ID, amount money.Money) error
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) Create(ctx context.Context, account Account) (Account, error) {
	return mock.OnCreate(ctx, account)
}

func (mock MockRepository) GetBalance(ctx context.Context, id ID) (int, error) {
	return mock.OnGetBalance(ctx, id)
}

func (mock MockRepository) Fetch(ctx context.Context) ([]Account, error) {
	return mock.OnFetch(ctx)
}

func (mock MockRepository) GetByID(ctx context.Context, id ID) (Account, error) {
	return mock.OnGetByID(ctx, id)
}

func (mock MockRepository) GetByCpf(ctx context.Context, cpf cpf.Cpf) (Account, error) {
	return mock.OnGetByCpf(ctx, cpf)
}

func (mock MockRepository) CreditAccount(ctx context.Context, id ID, amount money.Money) error {
	return mock.OnCreditAccount(ctx, id, amount)
}

func (mock MockRepository) DebitAccount(ctx context.Context, id ID, amount money.Money) error {
	return mock.OnDebitAccount(ctx, id, amount)
}
