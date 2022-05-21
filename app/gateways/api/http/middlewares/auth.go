package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/vos/token"
	"github.com/jpgsaraceni/suricate-bank/config"
)

type originIDKey struct{}

func Authorize(cfg config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := responses.NewResponse(w)

			authHeader := r.Header.Get("Authorization")
			requestToken := strings.ReplaceAll(authHeader, "Bearer ", "")

			if requestToken == "" {
				response.Unauthorized(responses.ErrMissingAuthorizationHeader).SendJSON()

				return
			}

			originID, err := token.Verify(cfg, requestToken)
			if err != nil {
				response.Unauthorized(responses.ErrInvalidToken).SendJSON()

				return
			}

			ctx := WithOriginID(r.Context(), originID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OriginIDFromContext(ctx context.Context) (account.ID, bool) {
	originID, ok := ctx.Value(originIDKey{}).(account.ID)

	return originID, ok
}

func WithOriginID(ctx context.Context, originID account.ID) context.Context {
	return context.WithValue(ctx, originIDKey{}, originID)
}
