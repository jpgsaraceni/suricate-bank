package transferuc

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type Usecase struct {
	Repository transfer.Repository
	Crediter
	Debiter
}

type Debiter interface {
	Debit(id account.AccountId, amount money.Money) (account.Account, error)
}

type Crediter interface {
	Credit(id account.AccountId, amount money.Money) (account.Account, error)
}

var (
	ErrCreateTransfer           = errors.New("failed transfer")
	ErrCreateTransferRepository = errors.New("repository error")
)
