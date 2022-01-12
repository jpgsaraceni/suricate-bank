package accountsroute

import (
	"errors"
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

func ErrorResponse(w http.ResponseWriter, err error) error {
	switch {
	case errors.Is(err, account.ErrInvalidCpf):
		return responses.SendJSON(w, responses.BadRequest(err, responses.ErrInvalidCpf))
	case errors.Is(err, ErrMissingFields):
		return responses.SendJSON(w, responses.BadRequest(err, responses.ErrMissingFields))
	case errors.Is(err, ErrInvalidRequestPayload):
		return responses.SendJSON(w, responses.BadRequest(err, responses.ErrInvalidRequestPayload))
	case errors.Is(err, ErrLengthCpf):
		return responses.SendJSON(w, responses.BadRequest(err, responses.ErrLengthCpf))
	case errors.Is(err, ErrShortName):
		return responses.SendJSON(w, responses.BadRequest(err, responses.ErrShortName))
	case errors.Is(err, ErrShortSecret):
		return responses.SendJSON(w, responses.BadRequest(err, responses.ErrShortSecret))
	default:
		return responses.SendJSON(w, responses.InternalServerError(err))
	}
}
