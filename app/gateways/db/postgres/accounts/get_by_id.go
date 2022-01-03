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

	var accountReturned queryReturn

	row := r.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&accountReturned.id,
		&accountReturned.name,
		&accountReturned.cpf,
		&accountReturned.secret,
		&accountReturned.balance,
		&accountReturned.createdAt,
	)

	if err != nil {

		return account.Account{}, fmt.Errorf("%w: %s", errQuery, err.Error())
	}

	parsedAccount, err := accountReturned.parse()

	if err != nil {

		return account.Account{}, fmt.Errorf("%w: %s", errParse, err.Error())
	}

	return parsedAccount, nil
}
