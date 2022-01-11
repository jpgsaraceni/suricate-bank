package accountuc

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (uc usecase) Credit(ctx context.Context, id account.AccountId, amount money.Money) error {
	account, err := uc.GetById(ctx, id)

	if err != nil {

		return fmt.Errorf("failed to get account by id: %w", err)
	}

	err = account.Balance.Add(amount.Cents())

	if err != nil {

		return fmt.Errorf("%w: %s", ErrAmount, err.Error())
	}

	err = uc.repository.CreditAccount(ctx, account.Id, amount)

	if err != nil {

		return fmt.Errorf("%w: %s", ErrCreditAccount, err.Error())
	}

	return nil
}
