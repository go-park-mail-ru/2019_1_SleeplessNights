package middleware

import "net/http"

const (
	DomainsCORS     = "http://localhost:8000"
	MethodsCORS     = "GET, POST, PATCH, DELETE, OPTIONS"
	CredentialsCORS = "true"
	//TODO FIX CORS HEADERS
	HeadersCORS = "X-Requested-With, Content-type, User-Agent, Cache-Control, Cookie, Origin, Accept-Encoding, Connection, Host, Upgrade-Insecure-Requests, User-Agent, Referer, Access-Control-Request-Method, Access-Control-Request-Headers"
)

func MiddlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", DomainsCORS)
		w.Header().Set("Access-Control-Allow-Credentials", CredentialsCORS)
		w.Header().Set("Access-Control-Allow-Methods", MethodsCORS)
		w.Header().Set("Access-Control-Allow-Headers", HeadersCORS)
		next.ServeHTTP(w, r)

	})
}
