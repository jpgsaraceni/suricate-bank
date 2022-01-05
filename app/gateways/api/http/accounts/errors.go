package accountsroute

import (
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

var (
	ErrMissingFields       = responses.ErrorPayload{Message: "missing fields: name, cpf and/or secret"}
	ErrLengthCpf           = responses.ErrorPayload{Message: "invalid cpf length"}
	ErrShortName           = responses.ErrorPayload{Message: "name too short"}
	ErrShortSecret         = responses.ErrorPayload{Message: "password too short"}
	ErrInternalServerError = responses.ErrorPayload{Message: "internal server error"}
)
