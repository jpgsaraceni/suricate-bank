package accountuc

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

type MockUsecase struct {
	OnCreate     func(ctx context.Context, accountInstance account.Account) (account.Account, error)
	OnGetBalance func(ctx context.Context, id account.ID) (int, error)
	OnFetch      func(ctx context.Context) ([]account.Account, error)
}

var _ Usecase = (*MockUsecase)(nil)

func (mock MockUsecase) Create(ctx context.Context, accountInstance account.Account) (account.Account, error) {
	return mock.OnCreate(ctx, accountInstance)
}

func (mock MockUsecase) GetBalance(ctx context.Context, id account.ID) (int, error) {
	return mock.OnGetBalance(ctx, id)
}

func (mock MockUsecase) Fetch(ctx context.Context) ([]account.Account, error) {
	return mock.OnFetch(ctx)
}
