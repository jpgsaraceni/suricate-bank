package accountuc

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

var errDebitAccountRepository = errors.New("repository error when debiting account")

func (uc Usecase) Debit(id account.AccountId, amount money.Money) error {
	account, err := uc.GetById(id)

	if err != nil {

		return err // TODO translate this error
	}

	err = account.Balance.Subtract(amount.Cents())

	if err != nil {

		return err
	}

	err = uc.Repository.DebitAccount(&account, amount)

	if err != nil {

		return errDebitAccountRepository
	}

	return nil
}
