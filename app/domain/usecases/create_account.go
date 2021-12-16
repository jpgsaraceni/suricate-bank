package usecase

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

var (
	errCreate      = errors.New("failed to create account")
	errNameLength  = errors.New("invalid name length")
	errShortSecret = errors.New("invalid password length")
	errRepository  = errors.New("repository error")
)

const (
	minNameLength = 3
	maxNameLength = 100

	minPasswordLength = 6
)

type Usecase struct {
	Repository account.Repository
}

func (uc Usecase) Create(name, cpf, secret string) (account.Account, error) {
	var newAccount account.Account

	if len(name) < minNameLength || len(name) > maxNameLength {

		return newAccount, errNameLength
	}

	if len(secret) < minPasswordLength {

		return newAccount, errShortSecret
	}

	newAccount, err := account.NewAccount(name, cpf, secret)

	if err != nil {
		return newAccount, errCreate
	}

	err = uc.Repository.Create(&newAccount)

	if err != nil {
		return newAccount, errRepository
	}

	return newAccount, nil
}
