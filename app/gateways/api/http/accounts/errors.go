package accountsroute

import "errors"

var (
	ErrMissingFields = errors.New("missing fields: name, cpf and/or secret")
	ErrLengthCpf     = errors.New("invalid cpf length")
	ErrShortName     = errors.New("name too short")
	ErrShortSecret   = errors.New("password too short")
)
