package transferuc

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (uc Usecase) Create(amount money.Money, originId, destinationId account.AccountId) (transfer.Transfer, error) {

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
