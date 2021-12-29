package accountuc

import (
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (uc Usecase) Debit(id account.AccountId, amount money.Money) error {
	account, err := uc.GetById(id)

	if err != nil {

		return fmt.Errorf("failed to get account by id: %w", err)
	}

	err = account.Balance.Subtract(amount.Cents())

	if err != nil {

		return fmt.Errorf("failed to debit amount: %w", err)
	}

	err = uc.Repository.DebitAccount(&account, amount)

	if err != nil {

		return fmt.Errorf("%w: %s", ErrDebitAccount, err.Error())
	}

	return nil
}
