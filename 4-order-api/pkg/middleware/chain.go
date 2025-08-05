package middleware

import (
	"go/order-api/configs"
	"net/http"
)

type Middleware func(http.Handler, *configs.Config) http.Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler, cfg *configs.Config) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next, cfg)
		}
		return next
	}
}
