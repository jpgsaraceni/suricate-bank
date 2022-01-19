package accountuc

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

// usecase calls Repository to be used in all methods of this package.
type usecase struct {
	repository account.Repository
}

type Usecase interface {
	Create(ctx context.Context, name, cpf, secret string) (account.Account, error)
	GetBalance(ctx context.Context, id account.AccountId) (int, error)
	Fetch(ctx context.Context) ([]account.Account, error)
}

func NewUsecase(r account.Repository) Usecase {
	return &usecase{
		repository: r,
	}
}
