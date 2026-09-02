package auth

import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

func LoginSubmitHandler(logger *slog.Logger, authSvc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			logger.Error("error while parsing form", "err", err)
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		ip := getClientIP(r)

		email := r.FormValue("email")
		password := r.FormValue("password")

		token, err := authSvc.LoginUserWithIPLimit(r.Context(), ip, email, password)
		if err != nil {
			if errors.Is(err, ErrTooManyAttempts) {
				http.Error(w, "too many attempts", http.StatusTooManyRequests)
				return
			}

			logger.Error("error occured during login user", "err", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			HttpOnly: r.TLS != nil,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   60 * 60 * 24,
		})

		http.Redirect(w, r, "/admin/panel", http.StatusSeeOther)
	}
}

func LogoutHandler(authSvc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("auth_token")
		authSvc.LogoutUser(r.Context(), c.Value)

		// delete cookie
		cookie := &http.Cookie{
			Name:     "auth_token",
			Value:    "",
			MaxAge:   -1,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}
}

func getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		ip := strings.TrimSpace(ips[0])
		if ip != "" {
			return ip
		}
	}

	if ip := strings.TrimSpace(r.Header.Get("X-Real-IP")); ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
