package middleware

import (
	"net/http"
	"strings"
)

func IsAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		_ = strings.TrimPrefix(header, "Bearer")
		//fmt.Println(token)
		next.ServeHTTP(w, r)
	})
}
