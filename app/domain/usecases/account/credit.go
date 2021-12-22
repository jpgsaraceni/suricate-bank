package accountuc

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (uc Usecase) Credit(id account.AccountId, amount money.Money) error {
	account, err := uc.GetById(id)

	if err != nil {

		return err // TODO translate error
	}

	err = account.Balance.Add(amount.Cents())

	if err != nil {

		return err
	}

	err = uc.Repository.CreditAccount(&account, amount)

	if err != nil {

		return err
	}

	return nil
}
