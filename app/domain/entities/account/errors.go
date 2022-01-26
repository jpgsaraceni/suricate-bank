package account

import "errors"

var (
	ErrInvalidCpf        = errors.New("invalid cpf")
	ErrEmptyName         = errors.New("empty name")
	ErrEmptySecret       = errors.New("empty secret")
	ErrIdNotFound        = errors.New("account id not found")
	ErrDuplicateCpf      = errors.New("cpf already assigned to existent account")
	ErrInsufficientFunds = errors.New("debit amount is greater than balance")
)
