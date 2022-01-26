package accountuc

import (
	"context"
	"errors"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (uc usecase) Credit(ctx context.Context, id account.AccountId, amount money.Money) error {
	if amount.Cents() == 0 {

		return ErrAmount
	}

	err := uc.repository.CreditAccount(ctx, id, amount)

	if err != nil {
		if errors.Is(err, account.ErrIdNotFound) {

			return err
		}

		return fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}

	return nil
}
