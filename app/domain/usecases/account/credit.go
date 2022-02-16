package accountuc

import (
	"context"
	"errors"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (uc usecase) Credit(ctx context.Context, id account.ID, amount money.Money) error {
	if amount.Cents() == 0 {
		return ErrAmount
	}

	err := uc.repository.CreditAccount(ctx, id, amount)
	if err != nil {
		if errors.Is(err, account.ErrIDNotFound) {
			return err
		}

		return fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}

	return nil
}
