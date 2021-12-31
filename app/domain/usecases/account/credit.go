package accountuc

import (
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (uc Usecase) Credit(id account.AccountId, amount money.Money) error {
	account, err := uc.GetById(id)

	if err != nil {

		return fmt.Errorf("failed to get account by id: %w", err)
	}

	err = account.Balance.Add(amount.Cents())

	if err != nil {

		return fmt.Errorf("%w: %s", ErrAmount, err.Error())
	}

	err = uc.Repository.CreditAccount(&account, amount)

	if err != nil {

		return fmt.Errorf("%w: %s", ErrCreditAccount, err.Error())
	}

	return nil
}
