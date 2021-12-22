package transferuc

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

// tudo bem ter esse erro repetido na entidade e na usecase?
var errSameAccounts = errors.New("origin and destination must be different accounts")

func (uc Usecase) Create(amount money.Money, originId, destinationId account.AccountId) (transfer.Transfer, error) {

	if originId == destinationId {

		return transfer.Transfer{}, errSameAccounts
	}

	err := uc.Debiter.Debit(originId, amount) // deveria não retornar a conta?

	if err != nil {

		return transfer.Transfer{}, err // TODO translate this error
	}

	// doesn't exist, negative or zero amount
	err = uc.Crediter.Credit(destinationId, amount)

	if err != nil {
		uc.Crediter.Credit(originId, amount) // deveria tratar o erro que pode retornar aqui?

		return transfer.Transfer{}, err
	}

	newTransfer, err := transfer.NewTransfer(amount, originId, destinationId)

	if err != nil {
		return transfer.Transfer{}, ErrCreateTransfer
	}

	err = uc.Repository.Create(&newTransfer)

	if err != nil {

		return transfer.Transfer{}, ErrCreateTransferRepository
	}

	return newTransfer, nil
}
