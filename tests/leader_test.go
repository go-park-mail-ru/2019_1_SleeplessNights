package tests

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLeadersHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/leaders", nil)
	if err != nil {
		t.Error(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.LeadersHandler)

	handler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf(WrongStatus+": got %v want %v",
			status, http.StatusOK)
	}

	expected := `{}` //TODO expected
	if resp.Body.String() != expected {
		t.Errorf(UnexpectedBody+": got %v want %v",
			resp.Body.String(), expected)
	}
}
