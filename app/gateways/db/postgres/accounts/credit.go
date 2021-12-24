package accountspg

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (r Repository) CreditAccount(account *account.Account, amount money.Money) error {
	const query = `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2;
	`

	err := r.pool.QueryRow(context.TODO(), query, amount.Cents(), account.Id)

	if err != nil {

		return errQuery
	}

	return nil
}
