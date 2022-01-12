package accountsroute

import "errors"

var (
	ErrLengthCpf             = errors.New("invalid cpf length")
	ErrShortName             = errors.New("name too short")
	ErrShortSecret           = errors.New("secret too short")
	ErrMissingFields         = errors.New("missing fields")
	ErrInvalidRequestPayload = errors.New("invalid request payload")
)
