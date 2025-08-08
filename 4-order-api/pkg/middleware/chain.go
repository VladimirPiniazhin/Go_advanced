package middleware

import (
	"go/order-api/configs"
	"net/http"
)

// Middleware с конфигурацией
type ConfigMiddleware func(http.Handler, *configs.Config) http.Handler

// Простой middleware без конфигурации
type SimpleMiddleware func(http.Handler) http.Handler

// Универсальный тип middleware
type Middleware interface {
	Apply(http.Handler, *configs.Config) http.Handler
}

// Адаптер для ConfigMiddleware
type configMiddleware struct {
	fn ConfigMiddleware
}

func (cm configMiddleware) Apply(next http.Handler, cfg *configs.Config) http.Handler {
	return cm.fn(next, cfg)
}

// Адаптер для SimpleMiddleware
type simpleMiddleware struct {
	fn SimpleMiddleware
}

func (sm simpleMiddleware) Apply(next http.Handler, cfg *configs.Config) http.Handler {
	return sm.fn(next)
}

// Конструкторы адаптеров
func WithConfig(fn ConfigMiddleware) Middleware {
	return configMiddleware{fn: fn}
}

func Simple(fn SimpleMiddleware) Middleware {
	return simpleMiddleware{fn: fn}
}

// Цепочка middleware
func Chain(middlewares ...Middleware) ConfigMiddleware {
	return func(next http.Handler, cfg *configs.Config) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i].Apply(next, cfg)
		}
		return next
	}
}
