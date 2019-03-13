package tests

import (
	"bytes"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImgHandler(t *testing.T) {

	path := "/img/static/img/default_avatar.jpg"
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	http.HandlerFunc(handlers.ImgHandler).ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	if bytes.Compare(resp.Body.Bytes(), expected) != 0 { //TODO
		t.Errorf("handler returned unexpected body")
	}
}
