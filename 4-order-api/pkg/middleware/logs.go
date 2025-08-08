package middleware

import (
	"go/order-api/configs"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Простая версия Logging без конфига
func LoggingSimple(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, request)
		log.WithFields(log.Fields{
			"statuscode": wrapper.StatusCode,
			"method":     request.Method,
			"path":       request.URL.Path,
			"execution":  time.Since(start).String(),
		}).Info("Handled request")
	})
}

// Deprecated: Use LoggingSimple instead
func Logging(next http.Handler, _ *configs.Config) http.Handler {
	return LoggingSimple(next)
}

// func Logging(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
// 		start := time.Now()
// 		wrapper := &WrapperWriter{
// 			ResponseWriter: w,
// 			StatusCode:     http.StatusOK,
// 		}
// 		next.ServeHTTP(wrapper, request)
// 		fmt.Println(wrapper.StatusCode, request.Method, request.URL.Path, time.Since(start))
// 	})
// }
