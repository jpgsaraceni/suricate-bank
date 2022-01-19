package accountuc

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

type MockUsecase struct {
	OnCreate     func(ctx context.Context, name, cpf, secret string) (account.Account, error)
	OnGetBalance func(ctx context.Context, id account.AccountId) (int, error)
	OnFetch      func(ctx context.Context) ([]account.Account, error)
}

var _ Usecase = (*MockUsecase)(nil)

func (mock MockUsecase) Create(ctx context.Context, name, cpf, secret string) (account.Account, error) {
	return mock.OnCreate(ctx, name, cpf, secret)
}

func (mock MockUsecase) GetBalance(ctx context.Context, id account.AccountId) (int, error) {
	return mock.OnGetBalance(ctx, id)
}

func (mock MockUsecase) Fetch(ctx context.Context) ([]account.Account, error) {
	return mock.OnFetch(ctx)
}
