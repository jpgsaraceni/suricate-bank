package transferuc

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

func (uc usecase) Create(ctx context.Context, transferInstance transfer.Transfer) (transfer.Transfer, error) {
	err := uc.Debiter.DebitAccount(ctx, transferInstance.AccountOriginID, transferInstance.Amount)
	if err != nil {
		return transfer.Transfer{}, fmt.Errorf("failed to debit origin account: %w", err)
	}

	err = uc.Crediter.CreditAccount(ctx, transferInstance.AccountDestinationID, transferInstance.Amount)

	if err != nil {
		rollbackCredit(ctx, uc, transferInstance)

		return transfer.Transfer{}, fmt.Errorf("failed to credit destination account: %w", err)
	}

	persistedTransfer, err := uc.Repository.Create(ctx, transferInstance)
	if err != nil {
		rollbackCredit(ctx, uc, transferInstance)
		rollbackDebit(ctx, uc, transferInstance)

		return transfer.Transfer{}, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}

	return persistedTransfer, nil
}

func rollbackDebit(ctx context.Context, uc usecase, transfer transfer.Transfer) {
	if err := uc.Debiter.DebitAccount(ctx, transfer.AccountOriginID, transfer.Amount); err != nil {
		log.Warn().Stack().Err(err).Msg("rollback failed")
	}
}

func rollbackCredit(ctx context.Context, uc usecase, transfer transfer.Transfer) {
	if err := uc.Crediter.CreditAccount(ctx, transfer.AccountDestinationID, transfer.Amount); err != nil {
		log.Warn().Stack().Err(err).Msg("rollback failed")
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
