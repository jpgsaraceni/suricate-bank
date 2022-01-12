package accountspg

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
)

func (r Repository) GetBalance(ctx context.Context, id account.AccountId) (int, error) {

	const query = `
		SELECT balance
		FROM accounts
		WHERE id = $1;
	`

	var balance int

	err := r.pool.QueryRow(ctx, query, id).Scan(&balance)

	log.Printf("type: %T error: %s", err, err)

	if errors.Is(err, pgx.ErrNoRows) {

		return 0, accountuc.ErrIdNotFound
	}

	if err != nil {

		return 0, ErrQuery
	}

	return balance, nil
}
