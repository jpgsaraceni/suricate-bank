package accountspg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (r Repository) Create(ctx context.Context, accountInstance account.Account) (account.Account, error) {
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
			($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			name,
			cpf,
			secret,
			balance,
			created_at
		;
	`

	var accountReturned account.Account

	err := r.pool.QueryRow(
		ctx,
		query,
		accountInstance.ID,
		accountInstance.Name,
		accountInstance.Cpf.Value(),
		accountInstance.Secret.Value(),
		accountInstance.Balance.Cents(),
		accountInstance.CreatedAt,
	).Scan(
		&accountReturned.ID,
		&accountReturned.Name,
		&accountReturned.Cpf,
		&accountReturned.Secret,
		&accountReturned.Balance,
		&accountReturned.CreatedAt,
	)

	const uniqueKeyViolationCode = "23505"
	const cpfConstraint = "accounts_cpf_key"

	var pgErr *pgconn.PgError

	if err != nil {
		if errors.As(err, &pgErr) {
			if pgErr.SQLState() == uniqueKeyViolationCode && pgErr.ConstraintName == cpfConstraint {
				return account.Account{}, account.ErrDuplicateCpf
			}
		}

		return account.Account{}, fmt.Errorf("%w: %s", ErrQuery, err.Error())
	}

	return accountReturned, nil
}
