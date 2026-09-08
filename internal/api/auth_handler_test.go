package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/MarcinBondaruk/fooder/internal/auth"
	"github.com/MarcinBondaruk/fooder/internal/auth/login_limiter"
	"github.com/MarcinBondaruk/fooder/internal/in_memory_db"
)

type fakeCredentialsRepo struct {
	getCredentialsByEmailFn func(ctx context.Context, email string) (*auth.UserCredentials, error)
}

func (f *fakeCredentialsRepo) GetCredentialsByEmail(ctx context.Context, email string) (*auth.UserCredentials, error) {
	return f.getCredentialsByEmailFn(ctx, email)
}

func newAuthService(credsFn func(context.Context, string) (*auth.UserCredentials, error), tokenStore map[string]struct{}, limiterLimit int) *auth.Service {
	credsRepo := &fakeCredentialsRepo{getCredentialsByEmailFn: credsFn}
	tokenRepo := in_memory_db.NewTokenRepository(tokenStore)
	limiter := login_limiter.NewLoginLimiter(limiterLimit, time.Hour)
	return auth.NewService(credsRepo, tokenRepo, limiter)
}

func TestLoginSubmitHandler(t *testing.T) {
	validCreds := func(_ context.Context, email string) (*auth.UserCredentials, error) {
		if email == "user@test.com" {
			return &auth.UserCredentials{ID: 1, Email: "user@test.com", HashedPassword: "secret"}, nil
		}
		return nil, errors.New("not found")
	}

	t.Run("successful login", func(t *testing.T) {
		tokenStore := make(map[string]struct{})
		svc := newAuthService(validCreds, tokenStore, 10)
		handler := LoginSubmitHandler(newDiscardLogger(), svc)

		form := url.Values{"email": {"user@test.com"}, "password": {"secret"}}
		req := httptest.NewRequest(http.MethodPost, "/admin/login-submit", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusSeeOther)
		}

		loc := rec.Header().Get("Location")
		if loc != "/admin/panel" {
			t.Errorf("Location = %q, want %q", loc, "/admin/panel")
		}

		cookies := rec.Result().Cookies()
		var authCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "auth_token" {
				authCookie = c
				break
			}
		}
		if authCookie == nil {
			t.Fatal("expected auth_token cookie to be set")
		}
		if authCookie.Value == "" {
			t.Error("expected non-empty auth_token cookie value")
		}
		if authCookie.Path != "/" {
			t.Errorf("cookie Path = %q, want %q", authCookie.Path, "/")
		}
		if authCookie.MaxAge != 86400 {
			t.Errorf("cookie MaxAge = %d, want 86400", authCookie.MaxAge)
		}

		if len(tokenStore) != 1 {
			t.Errorf("token store size = %d, want 1", len(tokenStore))
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		svc := newAuthService(validCreds, make(map[string]struct{}), 10)
		handler := LoginSubmitHandler(newDiscardLogger(), svc)

		form := url.Values{"email": {"user@test.com"}, "password": {"wrong"}}
		req := httptest.NewRequest(http.MethodPost, "/admin/login-submit", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		svc := newAuthService(validCreds, make(map[string]struct{}), 10)
		handler := LoginSubmitHandler(newDiscardLogger(), svc)

		form := url.Values{"email": {"nobody@test.com"}, "password": {"pass"}}
		req := httptest.NewRequest(http.MethodPost, "/admin/login-submit", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("rate limited", func(t *testing.T) {
		svc := newAuthService(validCreds, make(map[string]struct{}), 1)
		handler := LoginSubmitHandler(newDiscardLogger(), svc)

		form := url.Values{"email": {"nobody@test.com"}, "password": {"wrong"}}

		// First request: auth fails, registers one attempt (limit=1, now full)
		req1 := httptest.NewRequest(http.MethodPost, "/admin/login-submit", strings.NewReader(form.Encode()))
		req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec1 := httptest.NewRecorder()
		handler.ServeHTTP(rec1, req1)

		if rec1.Code != http.StatusUnauthorized {
			t.Errorf("first request: status = %d, want %d", rec1.Code, http.StatusUnauthorized)
		}

		// Second request: limiter is full, should return 429
		req2 := httptest.NewRequest(http.MethodPost, "/admin/login-submit", strings.NewReader(form.Encode()))
		req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec2 := httptest.NewRecorder()
		handler.ServeHTTP(rec2, req2)

		if rec2.Code != http.StatusTooManyRequests {
			t.Errorf("second request: status = %d, want %d", rec2.Code, http.StatusTooManyRequests)
		}
	})

	t.Run("empty form fields", func(t *testing.T) {
		svc := newAuthService(validCreds, make(map[string]struct{}), 10)
		handler := LoginSubmitHandler(newDiscardLogger(), svc)

		form := url.Values{"email": {""}, "password": {""}}
		req := httptest.NewRequest(http.MethodPost, "/admin/login-submit", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})
}

func TestLogoutHandler(t *testing.T) {
	t.Run("successful logout", func(t *testing.T) {
		tokenStore := map[string]struct{}{"tok123": {}}
		svc := newAuthService(nil, tokenStore, 10)
		handler := LogoutHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/admin/logout", nil)
		req.AddCookie(&http.Cookie{Name: "auth_token", Value: "tok123"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusSeeOther {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusSeeOther)
		}

		loc := rec.Header().Get("Location")
		if loc != "/admin/login" {
			t.Errorf("Location = %q, want %q", loc, "/admin/login")
		}

		cookies := rec.Result().Cookies()
		var authCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "auth_token" {
				authCookie = c
				break
			}
		}
		if authCookie == nil {
			t.Fatal("expected auth_token cookie in response")
		}
		if authCookie.MaxAge != -1 {
			t.Errorf("cookie MaxAge = %d, want -1 (deletion)", authCookie.MaxAge)
		}

		if _, exists := tokenStore["tok123"]; exists {
			t.Error("expected token to be removed from store")
		}
	})

	t.Run("no cookie panics", func(t *testing.T) {
		svc := newAuthService(nil, make(map[string]struct{}), 10)
		handler := LogoutHandler(svc)

		req := httptest.NewRequest(http.MethodGet, "/admin/logout", nil)
		rec := httptest.NewRecorder()

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when no auth_token cookie is present")
			}
		}()

		handler.ServeHTTP(rec, req)
	})
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		xff        string
		xRealIP    string
		remoteAddr string
		want       string
	}{
		{
			name:       "X-Forwarded-For single IP",
			xff:        "1.2.3.4",
			remoteAddr: "5.6.7.8:1234",
			want:       "1.2.3.4",
		},
		{
			name:       "X-Forwarded-For multiple IPs",
			xff:        "1.2.3.4, 10.0.0.1",
			remoteAddr: "5.6.7.8:1234",
			want:       "1.2.3.4",
		},
		{
			name:       "X-Real-IP only",
			xRealIP:    "9.8.7.6",
			remoteAddr: "5.6.7.8:1234",
			want:       "9.8.7.6",
		},
		{
			name:       "X-Forwarded-For takes precedence over X-Real-IP",
			xff:        "1.1.1.1",
			xRealIP:    "2.2.2.2",
			remoteAddr: "5.6.7.8:1234",
			want:       "1.1.1.1",
		},
		{
			name:       "RemoteAddr with port",
			remoteAddr: "192.168.1.1:9090",
			want:       "192.168.1.1",
		},
		{
			name:       "RemoteAddr without port",
			remoteAddr: "192.168.1.1",
			want:       "192.168.1.1",
		},
		{
			name:       "X-Forwarded-For empty falls through",
			xff:        "",
			remoteAddr: "10.0.0.1:80",
			want:       "10.0.0.1",
		},
		{
			name:       "whitespace in X-Forwarded-For",
			xff:        "  3.3.3.3 , 4.4.4.4",
			remoteAddr: "5.5.5.5:80",
			want:       "3.3.3.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.xff != "" {
				req.Header.Set("X-Forwarded-For", tt.xff)
			}
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}
			req.RemoteAddr = tt.remoteAddr

			got := getClientIP(req)
			if got != tt.want {
				t.Errorf("getClientIP() = %q, want %q", got, tt.want)
			}
		})
	}
}
