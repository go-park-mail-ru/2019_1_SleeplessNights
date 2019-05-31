package middleware_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	m "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"
	"strings"

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

	if resp.Header().Get("Access-Control-Allow-Origin") != strings.Join(config.GetStringSlice("shared.pkg.middleware.CORS.domains"), ",") {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			strings.Join(config.GetStringSlice("shared.pkg.middleware.CORS.domains"), ","), resp.Header().Get("Access-Control-Allow-Origin"))
	}

	if resp.Header().Get("Access-Control-Allow-Credentials") != config.GetString("shared.pkg.middleware.CORS.credentials") {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			config.GetString("shared.pkg.middleware.CORS.credentials"), resp.Header().Get("Access-Control-Allow-Credentials"))
	}

	if resp.Header().Get("Access-Control-Allow-Methods") != strings.Join(config.GetStringSlice("shared.pkg.middleware.CORS.methods"), ",") {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			strings.Join(config.GetStringSlice("shared.pkg.middleware.CORS.methods"), ","), resp.Header().Get("Access-Control-Allow-Methods"))
	}

	if resp.Header().Get("Access-Control-Allow-Headers") != strings.Join(config.GetStringSlice("shared.pkg.middleware.CORS.headers"), ",") {
		t.Errorf("Middleware get wrong header:\nwant: %v\ngot: %v",
			strings.Join(config.GetStringSlice("shared.pkg.middleware.CORS.headers"), ","), resp.Header().Get("Access-Control-Allow-Headers"))
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
