package accountspg

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (r Repository) Create(account *account.Account) error {

	const query = `
		INSERT INTO 
			accounts (
				id,
				name,
				cpf,
				secret,
				balance,
				created_at
			)
		VALUES
			($1, $2, $3, $4, $5, $6);
	`

	err := r.pool.QueryRow(
		context.TODO(),
		query,
		account.Id,
		account.Name,
		account.Cpf.Value(),
		account.Secret.Value(),
		account.Balance.Cents(),
		account.CreatedAt,
	).Scan()

	if err != nil {

		return fmt.Errorf("%w: %s", errQuery, err.Error())
	}

	return nil
}
