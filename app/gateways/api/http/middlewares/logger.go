package middlewares

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func RequestLogger(next http.Handler) http.Handler {
	loggerHandler := hlog.NewHandler(log.Logger)
	requestHandler := hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Str("ua", r.Header.Get("user-agent")).
			Str("ip", r.RemoteAddr).
			Str("referer", r.Referer()).
			Str("req_id", GetReqID(r.Context())).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	})

	return loggerHandler(requestHandler(next))
}
