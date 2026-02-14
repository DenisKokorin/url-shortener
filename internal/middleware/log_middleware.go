package logmiddleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func NewLogMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware"),
		)

		fn := func(w http.ResponseWriter, r *http.Request) {
			l := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t := time.Now()
			defer func() {
				l.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.String("duration", time.Since(t).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
