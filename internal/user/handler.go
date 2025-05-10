package user

import (
	"github.com/MarcinBondaruk/fooder/internal/auth/login_limiter"
	"net"
	"net/http"
	"strings"
)

func LoginSubmitHandler(ll *login_limiter.LoginLimiter, userSvc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		ip := getClientIP(r)

		email := r.FormValue("email")
		password := r.FormValue("password")

		token, err := userSvc.LoginUser(email, password)
		if err != nil {
			if !ll.Register(ip) {
				http.Error(w, "too many attempts", http.StatusUnauthorized)
				return
			}

			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ll.Release(ip)
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

func LogoutHandler(userSvc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("auth_token")

		userSvc.LogoutUser(c.Value)

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
