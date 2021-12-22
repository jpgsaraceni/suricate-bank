package accountuc

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (uc Usecase) Debit(id account.AccountId, amount money.Money) error {
	account, err := uc.GetById(id)

	if err != nil {

		return err
	}

	err = account.Balance.Subtract(amount.Cents())

	if err != nil {
		return err
	}

	return nil
}
