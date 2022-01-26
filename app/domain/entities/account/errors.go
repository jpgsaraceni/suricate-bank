package account

import "errors"

var (
	ErrInvalidCpf  = errors.New("invalid cpf")
	ErrEmptyName   = errors.New("empty name")
	ErrEmptySecret = errors.New("empty secret")
)
