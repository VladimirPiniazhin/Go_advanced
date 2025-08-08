package middleware

import (
	"go/order-api/configs"
	"net/http"
)

// WithAuth оборачивает хендлер в middleware авторизации
func WithAuth(handler http.HandlerFunc, config *configs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		IsAuthed(handler, config).ServeHTTP(w, r)
	}
}
