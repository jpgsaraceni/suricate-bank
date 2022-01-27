package responses

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

type Error struct {
	Err     error
	Payload ErrorPayload
}

var (
	// generic errors
	ErrInternalServerError = Error{
		Payload: ErrorPayload{Message: "Internal server error"},
	}

	// create account request errors
	ErrMissingFieldsAccountPayload = Error{
		Err:     errors.New("missing fields"),
		Payload: ErrorPayload{Message: "Missing fields: name, cpf and/or secret"},
	}
	ErrInvalidCreateAccountPayload = Error{
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
		Err:     account.ErrDuplicateCpf,
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

	// login errors
	ErrCredentials = Error{
		Err:     errors.New("cpf or secret do not match"),
		Payload: ErrorPayload{Message: "Incorrect CPF or password"},
	}
	ErrInvalidLoginPayload = Error{
		Err:     errors.New("invalid request payload"),
		Payload: ErrorPayload{Message: "Invalid payload. Expecting JSON containing keys cpf and secret"},
	}

	// create transfer error
	ErrMissingFieldsTransferPayload = Error{
		Err:     errors.New("missing fields"),
		Payload: ErrorPayload{Message: "Missing fields: destination account Id and/or amount"},
	}
	ErrMissingAuthorizationHeader = Error{
		Err:     errors.New("missing authorization header"),
		Payload: ErrorPayload{Message: "Missing header authorization bearer"},
	}
	ErrInvalidCreateTransferPayload = Error{
		Err:     errors.New("invalid request payload"),
		Payload: ErrorPayload{Message: "Invalid payload. Expecting JSON containing keys account_detination_id and amount"},
	}
	ErrInvalidDestinationId = Error{
		Err:     errors.New("failed to parse destination id request field to AccountIt"),
		Payload: ErrorPayload{Message: "Invalid account detination id"},
	}
	ErrInvalidAmount = Error{
		Err:     errors.New("amound not positive"),
		Payload: ErrorPayload{Message: "Amount to transfer must be grater than 0"},
	}
	ErrInsuficientFunds = Error{
		Err:     errors.New("insuficient balance in origin account"),
		Payload: ErrorPayload{Message: "Insuficient balance in account"},
	}
	ErrInvalidToken = Error{
		Err:     errors.New("failed to verify token"),
		Payload: ErrorPayload{Message: "Invalid authorization token"},
	}
	ErrSameAccounts = Error{
		Err:     errors.New("same origin and destination accounts"),
		Payload: ErrorPayload{Message: "Cannot transfer to same account"},
	}
)
