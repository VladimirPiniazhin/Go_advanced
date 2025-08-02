package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, request)
		fmt.Println(wrapper.StatusCode, request.Method, request.URL.Path, time.Since(start))
	})
}
