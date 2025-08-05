package middleware

import (
	"context"
	"fmt"
	"go/order-api/configs"
	"go/order-api/pkg/jwt"
	"net/http"
	"strings"
)

type key string

const (
	ContextPhoneKey key = "ContextPhoneKey"
)

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))

}

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			writeUnauthorized(w)
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		isValid, data := jwt.NewJWT(config.Jwt.Secret).Parse(token)
		if !isValid {
			writeUnauthorized(w)
			return
		}
		fmt.Println(data.Phone)
		ctx := context.WithValue(r.Context(), ContextPhoneKey, data.Phone)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
