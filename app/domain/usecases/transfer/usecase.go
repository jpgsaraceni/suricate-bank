package transferuc

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type usecase struct {
	Repository transfer.Repository
	Crediter
	Debiter
}

type Debiter interface {
	DebitAccount(ctx context.Context, id account.AccountId, amount money.Money) error
}

type Crediter interface {
	CreditAccount(ctx context.Context, id account.AccountId, amount money.Money) error
}

type Usecase interface {
	Create(ctx context.Context, transferInstance transfer.Transfer) (transfer.Transfer, error)
	Fetch(ctx context.Context) ([]transfer.Transfer, error)
}

func NewUsecase(r transfer.Repository, ar account.Repository) Usecase {
	return &usecase{
		Repository: r,
		Crediter:   ar,
		Debiter:    ar,
	}
}
