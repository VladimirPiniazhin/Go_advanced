package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(middlwares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlwares) - 1; i >= 0; i-- {
			next = middlwares[i](next)
		}
		return next
	}
}
