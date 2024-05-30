package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "development":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "production":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func LoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			duration := time.Since(start)

			logger.Info(
				"HTTP request",
				slog.Int("status", wrapped.status),
				slog.String("duration", duration.String()),
				slog.String("client_ip", r.RemoteAddr),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
		})
	}
}

// responseWriterWrapper обертка для http.ResponseWriter для захвата статуса ответа
type responseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
