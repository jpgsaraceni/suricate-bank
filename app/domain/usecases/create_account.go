package usecase

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

var (
	errCreate = errors.New("failed to create account")
	errName   = errors.New("invalid name length")
	errSecret = errors.New("invalid password length")
)

const (
	minNameLength = 3
	maxNameLength = 100

	minPasswordLength = 6
)

type Usecase struct {
	repository account.Repository
}

func (uc Usecase) Create(name, cpf, secret string) (account.Account, error) {
	var newAccount account.Account

	if len(name) < minNameLength || len(name) > maxNameLength {

		return newAccount, errName
	}

	if len(secret) < minPasswordLength {

		return newAccount, errSecret
	}

	newAccount, err := account.NewAccount(name, cpf, secret)

	if err != nil {
		return newAccount, errCreate
	}

	err = uc.repository.Create(&newAccount)

	if err != nil {
		return newAccount, errCreate
	}

	return newAccount, nil
}
