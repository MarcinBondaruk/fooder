package api

import (
	"log/slog"
	"net/http"
)

func CreateRecipe(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}
