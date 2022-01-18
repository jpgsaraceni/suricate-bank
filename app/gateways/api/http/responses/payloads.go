package responses

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
)

type Error struct {
	Err     error
	Payload ErrorPayload
}

var (
	// generic errors
	ErrInternalServerError = ErrorPayload{Message: "Internal server error"}

	// create account request errors
	ErrMissingFields = Error{
		Err:     errors.New("missing fields"),
		Payload: ErrorPayload{Message: "Missing fields: name, cpf and/or secret"},
	}
	ErrInvalidRequestPayload = Error{
		Err:     errors.New("invalid request payload"),
		Payload: ErrorPayload{Message: "Invalid payload. Expecting JSON containing keys name, cpf and secret"},
	}
	ErrLengthCpf = Error{
		Err:     errors.New("invalid cpf length"),
		Payload: ErrorPayload{Message: "CPF must contain 11 numeric digits or 14 digits including 3 symbols"},
	}
	ErrShortName = Error{
		Err:     errors.New("name too short"),
		Payload: ErrorPayload{Message: "Name must have at least 3 digits"},
	}
	ErrShortSecret = Error{
		Err:     errors.New("secret too short"),
		Payload: ErrorPayload{Message: "Secret must have at least 6 digits"},
	}
	ErrInvalidCpf = Error{
		Err:     account.ErrInvalidCpf,
		Payload: ErrorPayload{Message: "CPF is invalid"},
	}
	ErrCpfAlreadyExists = Error{
		Err:     accountspg.ErrCpfAlreadyExists,
		Payload: ErrorPayload{Message: "CPF already registered to an account"},
	}

	// get account balance request errors
	ErrInvalidPathParameter = Error{
		Err:     errors.New("invalid request url"),
		Payload: ErrorPayload{Message: "Invalid URL"},
	}
	ErrAccountNotFound = Error{
		Err:     errors.New("account not found"),
		Payload: ErrorPayload{Message: "Account not found"},
	}
)
