package accountuc

import (
	"context"
	"errors"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

const (
	minNameLength = 3
	maxNameLength = 100

	minPasswordLength = 6
)

func (uc usecase) Create(ctx context.Context, name, cpf, secret string) (account.Account, error) {
	if len(name) < minNameLength || len(name) > maxNameLength {

		return account.Account{}, ErrNameLength
	}

	if len(secret) < minPasswordLength {

		return account.Account{}, ErrShortSecret
	}

	newAccount, err := account.NewAccount(name, cpf, secret)

	if err != nil {

		return account.Account{}, fmt.Errorf("failed to create account instance: %w", err)
	}

	err = uc.repository.Create(ctx, &newAccount)

	if err != nil {
		if errors.Is(err, ErrDuplicateCpf) {

			return account.Account{}, err
		}

		return account.Account{}, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}

	return newAccount, nil
}
