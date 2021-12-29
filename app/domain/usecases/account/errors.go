package accountuc

import (
	"errors"
)

var (
	ErrNameLength        = errors.New("invalid name length")
	ErrShortSecret       = errors.New("invalid password length")
	ErrNoAccountsToFetch = errors.New("no accounts to fetch")
)
