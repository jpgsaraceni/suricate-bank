package accountuc

import (
	"errors"
)

var (
	ErrInsufficientFunds = errors.New("debit amount is greater than balance")
	ErrIdNotFound        = errors.New("account id not found")
	ErrDuplicateCpf      = errors.New("cpf already assigned to existent account")

	ErrNameLength  = errors.New("invalid name length")
	ErrShortSecret = errors.New("invalid password length")
	ErrAmount      = errors.New("invalid amount")

	ErrRepository = errors.New("accounts repository error")
)
