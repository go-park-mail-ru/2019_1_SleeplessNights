package tests

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOptionsHandler(t *testing.T) {
	req, err := http.NewRequest("PATCH", "/api/register", nil)
	if err != nil {
		t.Error(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.OptionsHandler)

	handler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusNoContent {
		t.Errorf(WrongStatus+": got %v want %v",
			status, http.StatusOK)
	}
}
