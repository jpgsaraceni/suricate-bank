package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func GetReqID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}

func RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}
