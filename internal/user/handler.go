package user

import "net/http"

func LoginSubmitHandler(userSvc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		token, err := userSvc.LoginUser(email, password)
		if err != nil {
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
