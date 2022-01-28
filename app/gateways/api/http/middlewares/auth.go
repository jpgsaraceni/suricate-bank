package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/vos/token"
)

type originIdKey struct{}

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

		ctx := WithOriginId(r.Context(), originId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OriginIdFromContext(ctx context.Context) (account.AccountId, bool) {
	originId, ok := ctx.Value(originIdKey{}).(account.AccountId)
	return originId, ok
}

func WithOriginId(ctx context.Context, originId account.AccountId) context.Context {
	return context.WithValue(ctx, originIdKey{}, originId)
}
