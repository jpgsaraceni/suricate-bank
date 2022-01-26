package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/vos/token"
)

type contextValueKey string

const ContextOriginId = contextValueKey("account_id")

func Authorize(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		response := responses.NewResponse(w)

		authHeader := r.Header.Get("Authorization")
		requestToken := strings.ReplaceAll(authHeader, "Bearer ", "")

		if requestToken == "" {
			response.Unauthorized(responses.ErrMissingAuthorizationHeader).SendJSON()

			return
		}

		originId, err := token.Verify(requestToken)

		if err != nil {
			response.Unauthorized(responses.ErrInvalidToken).SendJSON()

			return
		}

		ctx := context.WithValue(r.Context(), ContextOriginId, originId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
