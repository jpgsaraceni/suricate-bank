package responses

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

type Error struct {
	Err     error
	Payload Payload
}

var (
	AccountCreated         = Payload{Message: "Account successfully created"}
	ErrInternalServerError = Payload{Message: "Internal server error"}
)

var (
	// create account request errors
	ErrMissingFields = Error{
		Err:     errors.New("missing fields"),
		Payload: Payload{Message: "Missing fields: name, cpf and/or secret"},
	}
	ErrInvalidRequestPayload = Error{
		Err:     errors.New("invalid request payload"),
		Payload: Payload{Message: "Invalid payload. Expecting JSON containing keys name, cpf and secret"},
	}
	ErrLengthCpf = Error{
		Err:     errors.New("invalid cpf length"),
		Payload: Payload{Message: "CPF must contain 11 numeric digits or 14 digits including 3 symbols"},
	}
	ErrShortName = Error{
		Err:     errors.New("name too short"),
		Payload: Payload{Message: "Name must have at least 3 digits"},
	}
	ErrShortSecret = Error{
		Err:     errors.New("secret too short"),
		Payload: Payload{Message: "Secret must have at least 6 digits"},
	}
	ErrInvalidCpf = Error{
		Err:     account.ErrInvalidCpf,
		Payload: Payload{Message: "CPF is invalid"},
	}

	// get account balance request errors
	ErrInvalidPathParameter = Error{
		Err:     errors.New("invalid request url"),
		Payload: Payload{Message: "Invalid URL"},
	}
	ErrAccountNotFound = Error{
		Err:     errors.New("account not found"),
		Payload: Payload{Message: "Account not found"},
	}
)
