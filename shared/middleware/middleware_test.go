package middleware_test

import (
	m "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"

	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	DomainsCORS     = "http://localhost:8000"
	MethodsCORS     = "GET, POST, PATCH, DELETE, OPTIONS"
	CredentialsCORS = "true"
	//TODO FIX CORS HEADERS
	HeadersCORS = "X-Requested-With, Content-type, User-Agent, Cache-Control, Cookie, Origin, Accept-Encoding, Connection, Host, Upgrade-Insecure-Requests, User-Agent, Referer, Access-Control-Request-Method, Access-Control-Request-Headers"
)

func TestMiddlewareCORS(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	resp := httptest.NewRecorder()

	_router := mux.NewRouter()
	_router.Use(m.MiddlewareCORS)
	_router.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
	_router.ServeHTTP(resp, req)

	if resp.Header().Get("Access-Control-Allow-Origin") != DomainsCORS {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			DomainsCORS, resp.Header().Get("Access-Control-Allow-Origin"))
	}

	if resp.Header().Get("Access-Control-Allow-Credentials") != CredentialsCORS {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			CredentialsCORS, resp.Header().Get("Access-Control-Allow-Credentials"))
	}

	if resp.Header().Get("Access-Control-Allow-Methods") != MethodsCORS {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			MethodsCORS, resp.Header().Get("Access-Control-Allow-Methods"))
	}

	if resp.Header().Get("Access-Control-Allow-Headers") != HeadersCORS {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			HeadersCORS, resp.Header().Get("Access-Control-Allow-Headers"))
	}
}

func TestMiddlewareBasicHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	resp := httptest.NewRecorder()

	_router := mux.NewRouter()
	_router.Use(m.MiddlewareBasicHeaders)
	_router.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
	_router.ServeHTTP(resp, req)

	if resp.Header().Get("Content-type") != "application/json" {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			"application/json", resp.Header().Get("Content-type"))
	}
}
