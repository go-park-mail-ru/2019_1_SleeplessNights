package handlers_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOptionsHandlerSuccessful(t *testing.T) {

	req := httptest.NewRequest(http.MethodOptions, handlers.ApiRegister, nil)

	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.OptionsHandler).ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code:\n got %v\n want %v\n",
			status, http.StatusOK)
	}
}
