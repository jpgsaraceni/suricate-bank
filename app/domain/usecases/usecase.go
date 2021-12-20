package usecase

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

// Usecase calls Repository to be used in all methods of this package.
type Usecase struct {
	Repository account.Repository
}

var (
	ErrCreateAccount        = errors.New("failed to create account")
	ErrNameLength           = errors.New("invalid name length")
	ErrShortSecret          = errors.New("invalid password length")
	ErrRepository           = errors.New("repository error")
	ErrGetBalanceRepository = errors.New("failed to get balance")
)
