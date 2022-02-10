package middlewares

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency"
	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

type idempotencyKey struct{}

func Idempotency(s idempotency.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			requestIdempotencyKey := r.Header.Get("Idempotency-Key")

			if requestIdempotencyKey == "" { // client has not implemented idempotent requests
				next.ServeHTTP(w, r)

				return
			}

			ctx := WithIdempotency(r.Context(), requestIdempotencyKey)

			idempotentResponse, err := s.GetCachedResponse(ctx, requestIdempotencyKey)

			if err != nil {
				responses.NewResponse(w).InternalServerError(err).SendJSON()

				return
			}

			if idempotentResponse.Key != "" && idempotentResponse.ResponseBody == nil { // key exists in redis but
				// has not been populated yet
				responses.NewResponse(w).Processing().SendJSON()

				return
			}

			if reflect.DeepEqual(idempotentResponse, schema.CachedResponse{}) { // no cached response
				err := s.Lock(ctx, requestIdempotencyKey)

				if err != nil {
					responses.NewResponse(w).InternalServerError(err).SendJSON()

					return
				}

				hijackedWriter := NewResponseHijack(w)
				next.ServeHTTP(hijackedWriter, r)

				err = s.CacheResponse(
					ctx,
					schema.NewCachedResponse(
						requestIdempotencyKey,
						hijackedWriter.StatusCode,
						hijackedWriter.BodyBuffer.Bytes(),
					),
				)

				if err != nil {

					log.Printf("failed to cache response\nIdempotency-Key:%s\nError:%s", requestIdempotencyKey, err)
				}

				return
			}

			responses.NewResponse(w).SendCachedResponse(idempotentResponse)
		}
		return http.HandlerFunc(fn)
	}
}

func WithIdempotency(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, idempotencyKey{}, key)
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
