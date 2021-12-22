package transferuc // TODO refactor rollback

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

var (
	errSameAccounts = errors.New("origin and destination must be different accounts")
	errDebit        = errors.New("failed to debit from origin account")
	errCredit       = errors.New("failed to credit to destination account")
)

func (uc Usecase) Create(amount money.Money, originId, destinationId account.AccountId) (transfer.Transfer, error) {

	if originId == destinationId {

		return transfer.Transfer{}, errSameAccounts
	}

	err := uc.Debiter.Debit(originId, amount)

	if err != nil {

		return transfer.Transfer{}, errDebit
	}

	err = uc.Crediter.Credit(destinationId, amount)

	if err != nil {
		uc.Crediter.Credit(originId, amount)

		return transfer.Transfer{}, errCredit
	}

	newTransfer, err := transfer.NewTransfer(amount, originId, destinationId)

	if err != nil {
		uc.Debiter.Debit(destinationId, amount)
		uc.Crediter.Credit(originId, amount)

		return transfer.Transfer{}, ErrCreateTransfer
	}

	err = uc.Repository.Create(&newTransfer)

	if err != nil {
		uc.Debiter.Debit(destinationId, amount)
		uc.Crediter.Credit(originId, amount)

		return transfer.Transfer{}, ErrCreateTransferRepository
	}

	return newTransfer, nil
}
