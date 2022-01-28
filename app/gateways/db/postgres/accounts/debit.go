package accountspg

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (r Repository) DebitAccount(ctx context.Context, id account.AccountId, amount money.Money) error {
	const query = `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
		RETURNING id;
	`

	var updateId uuid.UUID

	err := r.pool.QueryRow(ctx, query, amount.Cents(), id).Scan(&updateId)

	const checkConstraintViolationCode = "23514"
	const balanceConstraint = "accounts_balance_check"

	var pgErr *pgconn.PgError

	if err != nil {
		if errors.As(err, &pgErr) {
			if pgErr.SQLState() == checkConstraintViolationCode && pgErr.ConstraintName == balanceConstraint {

				return account.ErrInsufficientFunds
			}
		}

		if errors.Is(err, pgx.ErrNoRows) {

			return account.ErrIdNotFound
		}

		return fmt.Errorf("%w: %s", ErrQuery, err.Error())
	}

	return nil
}
