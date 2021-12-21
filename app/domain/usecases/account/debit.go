package accountuc

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func Debit(account *account.Account, amount money.Money) error {
	err := account.Balance.Subtract(amount.Cents())

	if err != nil {
		return err
	}

	return nil
}
