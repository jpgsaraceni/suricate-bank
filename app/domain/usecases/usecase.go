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
	ErrCreateAccount           = errors.New("failed to create account")
	ErrNameLength              = errors.New("invalid name length")
	ErrShortSecret             = errors.New("invalid password length")
	ErrCreateAccountRepository = errors.New("repository error")
	ErrInvalidId               = errors.New("invalid id")
	ErrFetchAccounts           = errors.New("couldn't fetch accounts")
	ErrNoAccountsToFetch       = errors.New("no accounts to fetch")
)
