package router_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/router"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareCORS(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	resp := httptest.NewRecorder()

	_router := mux.NewRouter()
	_router.Use(router.MiddlewareCORS)
	_router.HandleFunc("/", func(http.ResponseWriter, *http.Request){})
	_router.ServeHTTP(resp, req)

	if resp.Header().Get("Access-Control-Allow-Origin") != router.DomainsCORS{
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			router.DomainsCORS, resp.Header().Get("Access-Control-Allow-Origin"))
	}

	if resp.Header().Get("Access-Control-Allow-Credentials") != router.CredentialsCORS{
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			router.CredentialsCORS, resp.Header().Get("Access-Control-Allow-Credentials"))
	}

	if resp.Header().Get("Access-Control-Allow-Methods") != router.MethodsCORS{
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			router.MethodsCORS, resp.Header().Get("Access-Control-Allow-Methods"))
	}

	if resp.Header().Get("Access-Control-Allow-Headers") != router.HeadersCORS{
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			router.HeadersCORS, resp.Header().Get("Access-Control-Allow-Headers"))
	}
}


func TestMiddlewareBasicHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	resp := httptest.NewRecorder()

	_router := mux.NewRouter()
	_router.Use(router.MiddlewareBasicHeaders)
	_router.HandleFunc("/", func(http.ResponseWriter, *http.Request){})
	_router.ServeHTTP(resp, req)

	if resp.Header().Get("Content-type") != "application/json"{
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			"application/json", resp.Header().Get("Content-type"))
	}
}