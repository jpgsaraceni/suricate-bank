package accountuc

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

const (
	minNameLength = 3
	maxNameLength = 100

	minPasswordLength = 6
)

func (uc Usecase) Create(name, cpf, secret string) (account.Account, error) {
	if len(name) < minNameLength || len(name) > maxNameLength {

		return account.Account{}, errNameLength
	}

	if len(secret) < minPasswordLength {

		return account.Account{}, errShortSecret
	}

	newAccount, err := account.NewAccount(name, cpf, secret)

	if err != nil {

		return account.Account{}, errCreateAccount
	}

	err = uc.Repository.Create(&newAccount)

	if err != nil {

		return account.Account{}, errCreateAccountRepository
	}

	return newAccount, nil
}
