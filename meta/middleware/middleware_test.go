package middleware_test

import (
	m "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/middleware"

	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareCORS(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	resp := httptest.NewRecorder()

	_router := mux.NewRouter()
	_router.Use(m.MiddlewareCORS)
	_router.HandleFunc("/", func(http.ResponseWriter, *http.Request) {})
	_router.ServeHTTP(resp, req)

	if resp.Header().Get("Access-Control-Allow-Origin") != m.DomainsCORS {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			m.DomainsCORS, resp.Header().Get("Access-Control-Allow-Origin"))
	}

	if resp.Header().Get("Access-Control-Allow-Credentials") != m.CredentialsCORS {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			m.CredentialsCORS, resp.Header().Get("Access-Control-Allow-Credentials"))
	}

	if resp.Header().Get("Access-Control-Allow-Methods") != m.MethodsCORS {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			m.MethodsCORS, resp.Header().Get("Access-Control-Allow-Methods"))
	}

	if resp.Header().Get("Access-Control-Allow-Headers") != m.HeadersCORS {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			m.HeadersCORS, resp.Header().Get("Access-Control-Allow-Headers"))
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
