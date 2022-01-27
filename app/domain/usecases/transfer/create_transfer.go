package transferuc

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

func (uc usecase) Create(ctx context.Context, transfer transfer.Transfer) error {

	err := uc.Debiter.Debit(ctx, transfer.AccountOriginId, transfer.Amount)

	if err != nil {

		return fmt.Errorf("failed to debit origin account: %w", err)
	}

	err = uc.Crediter.Credit(ctx, transfer.AccountDestinationId, transfer.Amount)

	if err != nil {
		rollback(ctx, uc, false, true, transfer)

		return fmt.Errorf("failed to credit destination account: %w", err)
	}

	err = uc.Repository.Create(ctx, &transfer)

	if err != nil {
		rollback(ctx, uc, true, true, transfer)

		return fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}

	return nil
}

func rollback(ctx context.Context, uc usecase, hasCredited, hasDebited bool, transfer transfer.Transfer) {
	if hasCredited {
		uc.Debiter.Debit(ctx, transfer.AccountOriginId, transfer.Amount)
	}

	if hasDebited {
		uc.Crediter.Credit(ctx, transfer.AccountDestinationId, transfer.Amount)
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
