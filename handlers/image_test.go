package handlers_test

import (
	"bytes"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImgHandler(t *testing.T) {

	path := "/img/default_avatar.jpg"

	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Add("Accept", "image/jpeg")

	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(path, handlers.ImgHandler)
	router.ServeHTTP(resp, req)

	//handlers.ImgHandler(resp, req)

	//http.HandlerFunc(handlers.ImgHandler).ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	if bytes.Compare(resp.Body.Bytes(), expected) != 0 {
		t.Errorf("handler returned unexpected body")
	}
}
