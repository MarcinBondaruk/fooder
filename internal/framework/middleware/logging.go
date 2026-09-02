package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type customResponse struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (cr *customResponse) WriteHeader(status int) {
	if cr.wroteHeader {
		return
	}

	cr.statusCode = status
	cr.wroteHeader = true
	cr.ResponseWriter.WriteHeader(status)
}

func (cr *customResponse) Write(b []byte) (int, error) {
	if !cr.wroteHeader {
		cr.WriteHeader(http.StatusOK)
	}

	return cr.ResponseWriter.Write(b)
}

func LoggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cr := &customResponse{
				statusCode:     200,
				ResponseWriter: w,
			}
			start := time.Now()
			next.ServeHTTP(cr, r)
			logger.Info("%s %s %d served in %v", r.Method, r.URL, cr.statusCode, time.Since(start))
		})
	}
}
