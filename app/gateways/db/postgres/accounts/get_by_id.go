package accountspg

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (r Repository) GetById(ctx context.Context, id account.AccountId) (account.Account, error) {
	const query = `
		SELECT 
			id,
			name,
			cpf,
			secret,
			balance,
			created_at
		FROM accounts
		WHERE id = $1;
	`

	var accountReturned account.Account

	row := r.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&accountReturned.Id,
		&accountReturned.Name,
		&accountReturned.Cpf,
		&accountReturned.Secret,
		&accountReturned.Balance,
		&accountReturned.CreatedAt,
	)

	if err != nil {

		return account.Account{}, fmt.Errorf("%w: %s", ErrQuery, err.Error())
	}

	return accountReturned, nil
}
