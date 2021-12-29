package accountuc

import (
	"errors"
)

var (
	ErrNameLength        = errors.New("invalid name length")
	ErrShortSecret       = errors.New("invalid password length")
	ErrNoAccountsToFetch = errors.New("no accounts to fetch")
	ErrCreateAccount     = errors.New("failed to save account")
	ErrCreditAccount     = errors.New("failed to save credit to account")
	ErrDebitAccount      = errors.New("failed to save debit to account")
	ErrFetchAccounts     = errors.New("failed to fetch accounts")
	ErrGetBalance        = errors.New("failed to get account balance")
	ErrGetAccount        = errors.New("failed to get account")
)
