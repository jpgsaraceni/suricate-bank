package responses

import (
	"encoding/json"
	"log"
	"net/http"
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

func Unmarshal(src []interface{}) (Response, error) {
	var status int
	var payload map[string]interface{}
	err := json.Unmarshal(src[0].([]byte), &status)

	if err != nil {

		return Response{}, err
	}

	err = json.Unmarshal(src[1].([]byte), &payload)

	if err != nil {

		return Response{}, err
	}

	return Response{Status: status, Payload: payload}, nil
}

// func (t *Response) RedisScan(src interface{}) (err error) {
// 	if t == nil {
// 		return fmt.Errorf("nil pointer")
// 	}
// 	switch src := src.(type) {
// 	case int:
// 		t.Status = src
// 	case map[string]interface{}:
// 		t.Payload = src
// 	default:
// 		err = fmt.Errorf("cannot convert from %T to %T", src, t)
// 	}
// 	return err
// }
