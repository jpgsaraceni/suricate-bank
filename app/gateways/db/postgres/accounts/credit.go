package accountspg

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (r Repository) CreditAccount(ctx context.Context, id account.AccountId, amount money.Money) error {
	const query = `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING id;
	`

	var updateId uuid.UUID

	err := r.pool.QueryRow(ctx, query, amount.Cents(), id).Scan(&updateId)

	if err != nil {

		return fmt.Errorf("%w: %s", ErrQuery, err.Error())
	}

	return nil
}
