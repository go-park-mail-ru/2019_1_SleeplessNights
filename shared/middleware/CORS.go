package middleware

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"net/http"
	"strings"
)

func MiddlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin",
			strings.Join(config.GetStringSlice("shared.pkg.middleware.CORS.domains"), ","))
		w.Header().Set("Access-Control-Allow-Credentials",
			config.GetString("shared.pkg.middleware.CORS.credentials"))
		w.Header().Set("Access-Control-Allow-Methods",
			strings.Join(config.GetStringSlice("shared.pkg.middleware.CORS.methods"), ","))
		w.Header().Set("Access-Control-Allow-Headers",
			strings.Join(config.GetStringSlice("shared.pkg.middleware.CORS.headers"), ","))
		next.ServeHTTP(w, r)

	})
}
