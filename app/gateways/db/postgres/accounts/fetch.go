package accountspg

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (r Repository) Fetch(ctx context.Context) ([]account.Account, error) {
	const query = `
		SELECT 
			id,
			name,
			cpf,
			secret,
			balance,
			created_at
		FROM accounts;
	`

	rows, err := r.pool.Query(ctx, query)

	if err != nil {

		return nil, fmt.Errorf("%w: %s", ErrQuery, err.Error())
	}

	defer rows.Close()
	var accountList []account.Account

	for rows.Next() {
		var accountReturned account.Account
		err := rows.Scan(
			&accountReturned.Id,
			&accountReturned.Name,
			&accountReturned.Cpf,
			&accountReturned.Secret,
			&accountReturned.Balance,
			&accountReturned.CreatedAt,
		)

		if err != nil {

			return nil, fmt.Errorf("%w: %s", ErrScanningRows, err.Error())
		}

		accountList = append(accountList, accountReturned)
	}

	return accountList, nil
}
