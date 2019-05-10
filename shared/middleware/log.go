package middleware

import "net/http"

func MiddlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Have some request on", r.URL)
		next.ServeHTTP(w, r)
	})
}
