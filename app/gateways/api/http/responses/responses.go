package responses

import (
	"encoding/json"
	"net/http"
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

func SendJSON(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response.Payload)
}
