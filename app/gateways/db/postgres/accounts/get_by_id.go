package accountspg

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (r Repository) GetById(id account.AccountId) (account.Account, error) {
	const query = `
		SELECT 
			id,
			name,
			cpf,
			secret,
			balance,
			created_at,
		FROM accounts
		WHERE id = $1;
	`

	var accountGot account.Account

	err := r.pool.QueryRow(context.TODO(), query, id).Scan(&accountGot) // TODO verify if types should be parsed

	if err != nil {

		return account.Account{}, errQuery
	}

	return accountGot, nil
}
