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
		var accountReturned queryReturn
		err := rows.Scan(
			&accountReturned.id,
			&accountReturned.name,
			&accountReturned.cpf,
			&accountReturned.secret,
			&accountReturned.balance,
			&accountReturned.createdAt,
		)

		if err != nil {

			return nil, fmt.Errorf("%w: %s", ErrScanningRows, err.Error())
		}

		parsedAccount, err := accountReturned.parse()

		if err != nil {

			return nil, fmt.Errorf("%w: %s", ErrParse, err.Error())
		}

		accountList = append(accountList, parsedAccount)
	}

	return accountList, nil
}
