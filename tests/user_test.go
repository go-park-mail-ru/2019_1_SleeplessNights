package tests

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/register", nil)
	if err != nil {
		t.Error(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.RegisterHandler)

	handler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf(WrongStatus+": got %v want %v; error: %s",
			status, http.StatusOK, resp.Body)
	}

	expected := `{}`
	if resp.Body.String() != expected {
		t.Errorf(UnexpectedBody+": got %v want %v",
			resp.Body.String(), expected)
	}
}
