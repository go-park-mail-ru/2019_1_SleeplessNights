package handlers_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLeadersHandlerSuccessful(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, handlers.ApiLeader, nil)
	qq := req.URL.Query()
	qq.Add("page", "1")
	req.URL.RawQuery = qq.Encode()

	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.LeadersHandler).ServeHTTP(resp, req)
	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: %v\nhandler can't write into responce \n",
			status)
	} else {
		if status := resp.Code; status != http.StatusOK {
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
				status, http.StatusOK)
		}
	}
}


func TestLeadersHandlerUnsuccessful_WithWrongValue(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, handlers.ApiLeader, nil)
	qq := req.URL.Query()
	qq.Add("page", "aa")
	req.URL.RawQuery = qq.Encode()

	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.LeadersHandler).ServeHTTP(resp, req)
	if status := resp.Code; status != http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
			status, http.StatusInternalServerError)
	}
}

func TestLeadersHandlerUnsuccessful_WithZeroValue(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, handlers.ApiLeader, nil)
	qq := req.URL.Query()
	qq.Add("page", "")
	req.URL.RawQuery = qq.Encode()

	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.LeadersHandler).ServeHTTP(resp, req)
	if status := resp.Code; status != http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
			status, http.StatusInternalServerError)
	}
}

