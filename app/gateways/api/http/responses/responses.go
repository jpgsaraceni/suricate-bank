package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

type Payload struct {
	Message string `json:"message"`
}

type Response struct {
	Status  int
	Error   error
	Payload interface{}
	Headers map[string]string // TODO: implement
	Writer  http.ResponseWriter
}

type ErrorPayload struct {
	Type    string      `json:"type" example:"srn:error:some_error"`
	Title   string      `json:"title,omitempty" example:"Message for some error"`
	Details interface{} `json:"details,omitempty" swaggertype:"object"`
}

func (r Response) BadRequest(err Error) Response {
	r.Status = http.StatusBadRequest
	r.Error = err.Err
	r.Payload = err.Payload
	return r
}

func (r Response) NotFound(err Error) Response {
	r.Status = http.StatusNotFound
	r.Error = err.Err
	r.Payload = err.Payload
	return r
}

func (r Response) InternalServerError(err error) Response {
	r.Status = http.StatusInternalServerError
	r.Error = err
	r.Payload = ErrInternalServerError
	return r
}

func (r Response) Ok(payload Payload) Response {
	r.Status = http.StatusOK
	r.Payload = payload
	return r
}

func (r Response) Created(payload Payload) Response {
	r.Status = http.StatusCreated
	r.Payload = payload
	return r
}

func (r Response) SendJSON() {
	r.Writer.Header().Set("Content-Type", "application/json")
	r.Writer.WriteHeader(r.Status)
	if err := json.NewEncoder(r.Writer).Encode(r.Payload); err != nil {
		log.Println(err) // TODO: fix after implementing log
	}
}
