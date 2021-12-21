package accountuc

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (uc Usecase) Debit(id account.AccountId, amount money.Money) (account.Account, error) {
	account, err := uc.GetById(id)

	if err != nil {

		return account, err
	}

	err = account.Balance.Subtract(amount.Cents())

	if err != nil {
		return account, err
	}

	return account, nil
}
