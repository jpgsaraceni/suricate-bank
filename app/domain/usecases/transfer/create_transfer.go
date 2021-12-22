package transferuc

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
	// errRollback     = errors.New("failed to rollback after transfer error")
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
		rollback(uc, false, true, originId, destinationId, amount)

		return transfer.Transfer{}, errCredit
	}

	newTransfer, err := transfer.NewTransfer(amount, originId, destinationId)

	if err != nil {
		rollback(uc, true, true, originId, destinationId, amount)

		return transfer.Transfer{}, ErrCreateTransfer
	}

	err = uc.Repository.Create(&newTransfer)

	if err != nil {
		rollback(uc, true, true, originId, destinationId, amount)

		return transfer.Transfer{}, ErrCreateTransferRepository
	}

	return newTransfer, nil
}

func rollback(uc Usecase, hasCredited, hasDebited bool, originId, destinationId account.AccountId, amount money.Money) {
	if hasCredited {
		uc.Debiter.Debit(destinationId, amount)
	}

	if hasDebited {
		uc.Crediter.Credit(originId, amount)
	}
}

// returning error:
// func rollback(uc Usecase, hasCredited, hasDebited bool, originId, destinationId account.AccountId, amount money.Money) error {
// 	if hasCredited {
// 		err := uc.Debiter.Debit(destinationId, amount)

// 		if err != nil {

// 			return errRollback
// 		}

// 	}
// 	if hasDebited {
// 		err := uc.Crediter.Credit(originId, amount)

// 		if err != nil {

// 			return errRollback
// 		}
// 	}

// 	return nil
// }
