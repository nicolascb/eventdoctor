package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging returns an http.Handler that logs every request with method, path,
// status code, duration, and remote address.
func Logging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.statusCode),
			slog.String("duration", fmt.Sprintf("%.2fms", float64(duration.Microseconds())/1000)),
			slog.String("remote_addr", r.RemoteAddr),
		}

		if rw.statusCode >= http.StatusInternalServerError {
			logger.LogAttrs(r.Context(), slog.LevelError, "request completed", attrs...)
		} else {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "request completed", attrs...)
		}
	})
}
