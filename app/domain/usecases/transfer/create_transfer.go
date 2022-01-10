package transferuc

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (uc Usecase) Create(ctx context.Context, amount money.Money, originId, destinationId account.AccountId) (transfer.Transfer, error) {

	if originId == destinationId {

		return transfer.Transfer{}, ErrSameAccounts
	}

	err := uc.Debiter.Debit(ctx, originId, amount)

	if err != nil {

		return transfer.Transfer{}, fmt.Errorf("%w: %s", ErrDebitOrigin, err.Error())
	}

	err = uc.Crediter.Credit(ctx, destinationId, amount)

	if err != nil {
		rollback(ctx, uc, false, true, originId, destinationId, amount)

		return transfer.Transfer{}, fmt.Errorf("%w: %s", ErrCreditDestination, err.Error())
	}

	newTransfer, err := transfer.NewTransfer(amount, originId, destinationId)

	if err != nil {
		rollback(ctx, uc, true, true, originId, destinationId, amount)

		return transfer.Transfer{}, fmt.Errorf("failed to create transfer instance: %w", err)
	}

	err = uc.Repository.Create(ctx, &newTransfer)

	if err != nil {
		rollback(ctx, uc, true, true, originId, destinationId, amount)

		return transfer.Transfer{}, fmt.Errorf("%w: %s", ErrCreateTransfer, err.Error())
	}

	return newTransfer, nil
}

func rollback(ctx context.Context, uc Usecase, hasCredited, hasDebited bool, originId, destinationId account.AccountId, amount money.Money) {
	if hasCredited {
		uc.Debiter.Debit(ctx, destinationId, amount)
	}

	if hasDebited {
		uc.Crediter.Credit(ctx, originId, amount)
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
