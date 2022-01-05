package accountsroute

import "github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"

var ErrMissingFields = responses.ErrorPayload{Message: "missing fields: name, cpf and/or secret"}
