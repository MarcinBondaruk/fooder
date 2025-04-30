package middleware

import (
	"github.com/MarcinBondaruk/fooder/internal/auth"
	"net/http"
)

func NewApiKeyAuthorization(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if apiKey != r.Header.Get("Authorization") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func NewCookieBasedAuthorization(authSvc *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("auth_token")

			if err != nil {
				http.Error(w, "no auth cookie", http.StatusUnauthorized)
				return
			}

			err = authSvc.VerifyToken(c.Value)
			if err != nil {
				http.Error(w, "invalid auth cookie", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
