package middleware

import (
	"log/slog"
	"net/http"
)

func PanicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				slog.Error("internal server error", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
