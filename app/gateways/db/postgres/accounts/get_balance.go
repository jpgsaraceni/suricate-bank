package accountspg

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (r Repository) GetBalance(id account.AccountId) (int, error) {

	const query = `
		SELECT balance
		FROM accounts
		WHERE id = $1;
	`

	var balance int

	err := r.pool.QueryRow(context.TODO(), query, id).Scan(&balance)

	if err != nil {

		return balance, errQuery
	}

	return balance, nil
}
