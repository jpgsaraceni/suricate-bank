package auth

import "errors"

var (
	ErrCredentials = errors.New("wrong password or cpf")
	ErrSignToken   = errors.New("failed to sign jwt")
	ErrCpfNotFound = errors.New("cpf not found")
)
