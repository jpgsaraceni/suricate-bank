package auth

import "errors"

var (
	ErrInexistentCpf = errors.New("inexistent cpf")
	ErrWrongPassword = errors.New("wrong password")
)
