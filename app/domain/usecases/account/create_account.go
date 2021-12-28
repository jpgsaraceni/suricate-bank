package accountuc

import (
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

const (
	minNameLength = 3
	maxNameLength = 100

	minPasswordLength = 6
)

func (uc Usecase) Create(name, cpf, secret string) (account.Account, error) {
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

	err = uc.Repository.Create(&newAccount)

	if err != nil {

		return account.Account{}, fmt.Errorf("repository failed to create account: %w", err)
	}

	return newAccount, nil
}
