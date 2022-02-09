package middlewares

import (
	"net/http"
	"reflect"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency"
	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

func Idempotency(next http.Handler, s idempotency.Service) http.Handler {
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
			next.ServeHTTP(w, r)
			// TODO: set cached response

			return
		}

		responses.NewResponse(w).UseCache(idempotentResponse).SendJSON()
	})
}
