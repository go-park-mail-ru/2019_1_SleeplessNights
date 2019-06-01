package middleware

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"net/http"
	"strings"
)

func MiddlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resMap := config.GetMapStringToInterface("shared.pkg.middleware.CORS.domains")
		var resSlice []string
		for _, value := range resMap {
			resSlice = append(resSlice, value.(string))
		}
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(resSlice, ","))
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		resMap = config.GetMapStringToInterface("shared.pkg.middleware.CORS.methods")
		resSlice = []string{}
		for _, value := range resMap {
			resSlice = append(resSlice, value.(string))
		}
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(resSlice, ","))

		resMap = config.GetMapStringToInterface("shared.pkg.middleware.CORS.headers")
		resSlice = []string{}
		for _, value := range resMap {
			resSlice = append(resSlice, value.(string))
		}
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(resSlice, ","))

		next.ServeHTTP(w, r)
	})
}
