package accountspg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

func (r Repository) GetByCpf(ctx context.Context, cpf cpf.Cpf) (account.Account, error) {
	const query = `
		SELECT 
			id,
			name,
			cpf,
			secret,
			balance,
			created_at
		FROM accounts
		WHERE cpf = $1;
	`

	var accountReturned account.Account

	row := r.pool.QueryRow(ctx, query, cpf.Value())

	err := row.Scan(
		&accountReturned.Id,
		&accountReturned.Name,
		&accountReturned.Cpf,
		&accountReturned.Secret,
		&accountReturned.Balance,
		&accountReturned.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {

		return account.Account{}, ErrCpfNotFound
	}

	if err != nil {

		return account.Account{}, fmt.Errorf("%w: %s", ErrQuery, err.Error())
	}

	return accountReturned, nil
}
