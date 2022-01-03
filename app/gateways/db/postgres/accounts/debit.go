package accountspg

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (r Repository) DebitAccount(ctx context.Context, id account.AccountId, amount money.Money) error {
	const query = `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2;
	`

	_, err := r.pool.Exec(ctx, query, amount.Cents(), id)

	if err != nil {

		return fmt.Errorf("%w: %s", errQuery, err.Error())
	}

	return nil
}
