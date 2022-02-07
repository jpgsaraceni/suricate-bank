package middlewares

import (
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency"
)

func Idempotency(next http.Handler, s idempotency.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idempotencyKey := r.Header.Get("Idempotency-Key")

		if idempotencyKey == "" {
			next.ServeHTTP(w, r)

			return
		}

		idempotentResponse, err := s.GetKeyValue(r.Context(), idempotencyKey)

		if err != nil { // TODO: refactor
			// TODO: set key-value
			next.ServeHTTP(w, r)

			return
		}

		idempotentResponse.SendJSON()
	})
}
