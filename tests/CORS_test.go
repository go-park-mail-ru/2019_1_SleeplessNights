package tests

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOptionsHandler(t *testing.T){
	req, err := http.NewRequest("PATCH", "/api/register", nil)
	if err != nil{
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.OptionsHandler)

	handler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{}` //TODO expected
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Body.String(), expected)
	}
}
