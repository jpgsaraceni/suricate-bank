package accountsroute

import "errors"

var ( // TODO: refatorar para criar payload e erro juntos
	ErrLengthCpf             = errors.New("invalid cpf length")
	ErrShortName             = errors.New("name too short")
	ErrShortSecret           = errors.New("secret too short")
	ErrMissingFields         = errors.New("missing fields")
	ErrInvalidRequestPayload = errors.New("invalid request payload")
	ErrInvalidPathParameter  = errors.New("invalid request url")
	ErrAccountNotFound       = errors.New("account not found")
)
