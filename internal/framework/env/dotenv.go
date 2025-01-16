package env

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Env Crate this with NewEnv.
type Env struct {
	allowedOrigins []string
	sqliteDsn      string
}

func NewEnv() (*Env, error) {
	allowedOrigins, err := initAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"))
	if err != nil {
		return nil, err
	}

	sqliteDsn := os.Getenv("SQLITE_DSN")
	if sqliteDsn == "" {
		return nil, errors.New("SQLITE_DSN is required")
	}

	return &Env{
		allowedOrigins: allowedOrigins,
		sqliteDsn:      sqliteDsn,
	}, nil
}

func (e *Env) AllowedOrigins() []string {
	return e.allowedOrigins
}

func (e *Env) SqliteDSN() string {
	return e.sqliteDsn
}

func initAllowedOrigins(allowedOrigins string) ([]string, error) {
	if allowedOrigins == "" {
		return nil, errors.New("env: ALLOWED_ORIGINS is empty")
	}

	origins := strings.Split(allowedOrigins, ";")
	out := make([]string, len(origins))

	for i, o := range origins {
		o = strings.TrimSpace(o)
		parsedUrl, err := url.Parse(o)
		if err != nil {
			return nil, errors.Join(errors.New("env: unparsable origin"), err)
		}

		if parsedUrl.Scheme != "https" && parsedUrl.Scheme != "http" {
			return nil, fmt.Errorf("env: unsupported scheme: %s", parsedUrl.Scheme)
		}

		if parsedUrl.Host == "" {
			return nil, fmt.Errorf("env: invalid host: %s", parsedUrl.Host)
		}

		out[i] = o
	}

	return out, nil
}
