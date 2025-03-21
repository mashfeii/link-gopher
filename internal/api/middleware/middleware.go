package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func SlogLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		slog.Group("request start",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		slog.Group("request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", lrw.statusCode),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
