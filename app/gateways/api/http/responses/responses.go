package responses

import (
	"encoding/json"
	"net/http"
)

type Payload struct {
	Message string `json:"message"`
}

type Response struct {
	Status  int
	Error   error
	Payload Payload
}

func BadRequest(err Error) Response {
	return Response{
		Status:  http.StatusBadRequest,
		Error:   err.Err,
		Payload: err.Payload,
	}
}

func NotFound(err Error) Response {
	return Response{
		Status:  http.StatusNotFound,
		Error:   err.Err,
		Payload: err.Payload,
	}
}

func InternalServerError(err error) Response {
	return Response{
		Status:  http.StatusInternalServerError,
		Error:   err,
		Payload: ErrInternalServerError,
	}
}

func Created(payload Payload) Response {
	return Response{
		Status:  http.StatusCreated,
		Payload: payload,
	}
}

func Ok(payload Payload) Response {
	return Response{
		Status:  http.StatusOK,
		Payload: payload,
	}
}

func SendJSON(w http.ResponseWriter, response Response) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	return json.NewEncoder(w).Encode(response.Payload)
}
