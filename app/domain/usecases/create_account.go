package usecase

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities"
)

var (
	errCreate = errors.New("failed to create account")
	errName   = errors.New("invalid name length")
	errCpf    = errors.New("invalid cpf length")
	errSecret = errors.New("invalid password length")
)

const (
	minNameLength = 3
	maxNameLength = 100

	rawCpfLength    = 11
	maskedCpfLength = 14

	minPasswordLength = 6
)

// CreateAccount checks if lengths of arguments are ok, then calls entities.NewAccount,
// and returns the created Account struct.
func CreateAccount(name, cpf, secret string) (entities.Account, error) {
	var newAccount entities.Account

	if len(name) < minNameLength || len(name) > maxNameLength {

		return newAccount, errName
	}

	if len(cpf) != rawCpfLength || len(cpf) != maskedCpfLength {

		return newAccount, errCpf
	}

	if len(secret) < minPasswordLength {

		return newAccount, errSecret
	}

	newAccount, err := entities.NewAccount(name, cpf, secret)

	if err != nil {
		return newAccount, errCreate
	}

	return newAccount, nil
}
