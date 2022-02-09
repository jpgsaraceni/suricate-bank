package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency"
	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

func Idempotency(s idempotency.Service, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idempotencyKey := r.Header.Get("Idempotency-Key")

		if idempotencyKey == "" { // client has not implemented idempotent requests
			next.ServeHTTP(w, r)

			return
		}

		idempotentResponse, err := s.GetCachedResponse(r.Context(), idempotencyKey)

		if err != nil {
			responses.NewResponse(w).InternalServerError(err).SendJSON()

			return
		}

		if reflect.DeepEqual(idempotentResponse, schema.CachedResponse{}) { // no cached response
			hijackedWriter := NewResponseHijack(w)
			next(hijackedWriter, r)

			err := s.CacheResponse(
				schema.NewCachedResponse(
					idempotencyKey,
					hijackedWriter.StatusCode,
					hijackedWriter.BodyBuffer.Bytes(),
				),
			)

			if err != nil {

				log.Printf("failed to cache response\nIdempotency-Key:%s\nError:%s", idempotencyKey, err)
			}

			return
		}

		responses.NewResponse(w).SendCachedResponse(idempotentResponse)
	})
}

// ResponseHijack writes a response to http.ResponseWriter and a copy to BodyBuffer and Status
// using io.MultiWriter()
type ResponseHijack struct {
	w          http.ResponseWriter
	multi      io.Writer
	BodyBuffer *bytes.Buffer
	StatusCode int
}

func NewResponseHijack(w http.ResponseWriter) *ResponseHijack {
	bodyBuff := &bytes.Buffer{}
	multi := io.MultiWriter(bodyBuff, w)
	return &ResponseHijack{
		w:          w,
		multi:      multi,
		BodyBuffer: bodyBuff,
	}
}

func (h *ResponseHijack) Header() http.Header {
	return h.w.Header()
}

func (h *ResponseHijack) Write(b []byte) (int, error) {
	return h.multi.Write(b)
}

func (h *ResponseHijack) WriteHeader(i int) {
	h.StatusCode = i
	h.w.WriteHeader(i)
}
