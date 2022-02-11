package responses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

type Response struct {
	Status  int
	Error   error
	Payload interface{}
	Headers map[string]string
	Writer  http.ResponseWriter
}

type ErrorPayload struct {
	Message string `json:"title,omitempty" example:"Message for some error"`
}

type ProcessingPayload struct {
	Message string `json:"title,omitempty" example:"Request is being processed"`
}

func NewResponse(w http.ResponseWriter) Response {
	return Response{Writer: w}
}

func (r Response) IsComplete() bool {
	return r.Status > 0
}

func (r Response) BadRequest(err Error) Response {
	r.Status = http.StatusBadRequest
	r.Error = err.Err
	r.Payload = err.Payload
	return r
}

func (r Response) Unauthorized(err Error) Response {
	r.Status = http.StatusUnauthorized
	r.Error = err.Err
	r.Payload = err.Payload
	return r
}

func (r Response) Forbidden(err Error) Response {
	r.Status = http.StatusForbidden
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

func (r Response) UnprocessableEntity(err Error) Response {
	r.Status = http.StatusUnprocessableEntity
	r.Error = err.Err
	r.Payload = err.Payload
	return r
}

func (r Response) InternalServerError(err error) Response {
	r.Status = http.StatusInternalServerError
	r.Error = err
	r.Payload = ErrInternalServerError.Payload
	return r
}

func (r Response) Ok(payload interface{}) Response {
	r.Status = http.StatusOK
	r.Payload = payload
	return r
}

func (r Response) Created(payload interface{}) Response {
	r.Status = http.StatusCreated
	r.Payload = payload
	return r
}

func (r Response) Processing() Response {
	r.Status = http.StatusBadRequest
	r.Payload = ProcessingPayload{Message: "Request is duplicate. Original request is being processed."}
	return r
}

func (r Response) SendJSON() {
	r.Writer.Header().Set("Content-Type", "application/json")
	r.Writer.WriteHeader(r.Status)
	for headerKey, headerValue := range r.Headers {
		r.Writer.Header().Set(headerKey, headerValue)
	}
	if err := json.NewEncoder(r.Writer).Encode(r.Payload); err != nil {
		log.Println(err) // TODO: fix after implementing log
	}
}

func (r Response) SendCachedResponse(cache schema.CachedResponse) {
	r.Writer.Header().Set("Content-Type", "application/json")
	r.Writer.WriteHeader(cache.ResponseStatus)
	if _, err := r.Writer.Write(cache.ResponseBody); err != nil {
		log.Println(err)
	}
}
