package accountspg

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
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
		if errors.Is(err, pgx.ErrNoRows) {

			return accountuc.ErrIdNotFound
		}

		return fmt.Errorf("%w: %s", ErrQuery, err.Error())
	}

	return nil
}
