package accountspg

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (r Repository) Fetch() ([]account.Account, error) {
	const query = `
		SELECT *
		FROM accounts;
	`

	rows, err := r.pool.Query(context.TODO(), query, nil)

	if err != nil {

		return nil, errQuery
	}

	defer rows.Close()
	var accountList []account.Account

	for rows.Next() {
		var account account.Account
		err := rows.Scan( // TODO verify if types should be parsed
			&account.Id,
			&account.Name,
			&account.Cpf,
			&account.Secret,
			&account.Balance,
			&account.CreatedAt,
		)

		if err != nil {

			return nil, errScanningRows
		}

		accountList = append(accountList, account)
	}

	return accountList, nil
}
