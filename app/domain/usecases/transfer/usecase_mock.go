package transferuc

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type MockUsecase struct {
	OnCreate func(ctx context.Context, transferInstance transfer.Transfer) (transfer.Transfer, error)
	OnFetch  func(ctx context.Context) ([]transfer.Transfer, error)
}

var _ Usecase = (*MockUsecase)(nil)

func (mock MockUsecase) Create(ctx context.Context, transferInstance transfer.Transfer) (transfer.Transfer, error) {
	return mock.OnCreate(ctx, transferInstance)
}

func (mock MockUsecase) Fetch(ctx context.Context) ([]transfer.Transfer, error) {
	return mock.OnFetch(ctx)
}

type MockDebiter struct {
	OnDebitAccount func(ctx context.Context, id account.ID, amount money.Money) error
}

var _ Debiter = (*MockDebiter)(nil)

func (mock MockDebiter) DebitAccount(ctx context.Context, id account.ID, amount money.Money) error {
	return mock.OnDebitAccount(ctx, id, amount)
}

type MockCrediter struct {
	OnCreditAccount func(ctx context.Context, id account.ID, amount money.Money) error
}

var _ Crediter = (*MockCrediter)(nil)

func (mock MockCrediter) CreditAccount(ctx context.Context, id account.ID, amount money.Money) error {
	return mock.OnCreditAccount(ctx, id, amount)
}
