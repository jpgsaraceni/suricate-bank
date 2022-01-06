package responses

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

type Response struct {
	Status  int
	Error   error
	Payload interface{}
}

type ErrorPayload struct {
	Message string
}

func BadRequest(err error, errorPayload ErrorPayload) Response {
	return Response{
		Status:  http.StatusBadRequest,
		Error:   err,
		Payload: errorPayload,
	}
}

func InternalServerError(err error) Response {
	return Response{
		Status:  http.StatusInternalServerError,
		Error:   err,
		Payload: ErrInternalServerError,
	}
}

func Created(payload string) Response {
	return Response{
		Status:  http.StatusCreated,
		Payload: payload,
	}
}

func SendJSON(w http.ResponseWriter, response Response) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	return json.NewEncoder(w).Encode(response.Payload)
}

func ErrorResponse(err error) Response {
	switch {
	case errors.Is(err, account.ErrInvalidCpf):
		return BadRequest(err, ErrInvalidCpf)
	default:
		return InternalServerError(err)
	}
}
