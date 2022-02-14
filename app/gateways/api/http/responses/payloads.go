package responses

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

type Error struct {
	Err     error
	Payload CustomPayload
}

var (
	// generic errors
	ErrInternalServerError = Error{
		Payload: CustomPayload{Message: "Internal server error"},
	}

	// create account request errors
	ErrMissingFieldsAccountPayload = Error{
		Err:     errors.New("missing fields"),
		Payload: CustomPayload{Message: "Missing fields: name, cpf and/or secret"},
	}
	ErrInvalidCreateAccountPayload = Error{
		Err:     errors.New("invalid request payload"),
		Payload: CustomPayload{Message: "Invalid payload. Expecting JSON containing keys name, cpf and secret"},
	}
	ErrLengthCpf = Error{
		Err:     errors.New("invalid cpf length"),
		Payload: CustomPayload{Message: "CPF must contain 11 numeric digits or 14 digits including 3 symbols"},
	}
	ErrLengthName = Error{
		Err:     errors.New("invalid name length"),
		Payload: CustomPayload{Message: "Name must have from 3 to 100 digits"},
	}
	ErrShortSecret = Error{
		Err:     errors.New("secret too short"),
		Payload: CustomPayload{Message: "Secret must have at least 6 digits"},
	}
	ErrInvalidCpf = Error{
		Err:     account.ErrInvalidCpf,
		Payload: CustomPayload{Message: "CPF is invalid"},
	}
	ErrCpfAlreadyExists = Error{
		Err:     account.ErrDuplicateCpf,
		Payload: CustomPayload{Message: "CPF already registered to an account"},
	}

	// get account balance request errors
	ErrInvalidPathParameter = Error{
		Err:     errors.New("invalid request url"),
		Payload: CustomPayload{Message: "Invalid URL"},
	}
	ErrAccountNotFound = Error{
		Err:     errors.New("account not found"),
		Payload: CustomPayload{Message: "Account not found"},
	}

	// login errors
	ErrCredentials = Error{
		Err:     errors.New("cpf or secret do not match"),
		Payload: CustomPayload{Message: "Incorrect CPF or password"},
	}
	ErrInvalidLoginPayload = Error{
		Err:     errors.New("invalid request payload"),
		Payload: CustomPayload{Message: "Invalid payload. Expecting JSON containing keys cpf and secret"},
	}

	// create transfer error
	ErrMissingFieldsTransferPayload = Error{
		Err:     errors.New("missing fields"),
		Payload: CustomPayload{Message: "Missing fields: destination account Id and/or amount"},
	}
	ErrMissingAuthorizationHeader = Error{
		Err:     errors.New("missing authorization header"),
		Payload: CustomPayload{Message: "Missing header authorization bearer"},
	}
	ErrInvalidCreateTransferPayload = Error{
		Err:     errors.New("invalid request payload"),
		Payload: CustomPayload{Message: "Invalid payload. Expecting JSON containing keys account_detination_id and amount"},
	}
	ErrInvalidDestinationId = Error{
		Err:     errors.New("failed to parse destination id request field to AccountIt"),
		Payload: CustomPayload{Message: "Invalid account detination id"},
	}
	ErrInvalidAmount = Error{
		Err:     errors.New("amound not positive"),
		Payload: CustomPayload{Message: "Amount to transfer must be grater than 0"},
	}
	ErrInsuficientFunds = Error{
		Err:     errors.New("insuficient balance in origin account"),
		Payload: CustomPayload{Message: "Insuficient balance in account"},
	}
	ErrInvalidToken = Error{
		Err:     errors.New("failed to verify token"),
		Payload: CustomPayload{Message: "Invalid authorization token"},
	}
	ErrSameAccounts = Error{
		Err:     errors.New("same origin and destination accounts"),
		Payload: CustomPayload{Message: "Cannot transfer to same account"},
	}
)
