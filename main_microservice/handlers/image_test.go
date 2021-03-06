package handlers_test

import (
	"bytes"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestImgHandlerSuccessful(t *testing.T) {

	img := "default_avatar.jpg"
	path := fmt.Sprintf("%s%s", handlers.Img, img)
	req := httptest.NewRequest(http.MethodGet, path, nil)

	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/img/{path}", handlers.ImgHandler)
	router.ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: %v\nhandler can't write into responce or can't read file\n",
			status)
	} else {
		if status := resp.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code:\ngot %v\nwant %v\n",
				status, http.StatusOK)
		}

		expected, err := ioutil.ReadFile(os.Getenv("BASEPATH") + handlers.AvatarPrefix + img)
		if err != nil {
			t.Error(err.Error())
			return
		}

		if bytes.Compare(resp.Body.Bytes(), expected) != 0 {
			t.Errorf("handler returned unexpected body\n")
		}
	}
}

func TestImgHandlerUnsuccessful_WrongImagePath(t *testing.T) {

	img := "WRONG_default_avatar.jpg"
	path := fmt.Sprintf("%s%s", handlers.Img, img)
	req := httptest.NewRequest(http.MethodGet, path, nil)

	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/img/{path}", handlers.ImgHandler)
	router.ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: %v\nhandler can't write into responce or can't read file\n",
			status)
	} else {
		if status := resp.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code:\ngot %v\nwant %v\n",
				status, http.StatusNotFound)
		}
	}
}

func TestImgHandlerUnsuccessful_WithoutPath(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/img/", nil)

	resp := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/img/", handlers.ImgHandler)
	router.ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: %v\nhandler can't write into responce or can't read file\n",
			status)
	} else {
		if status := resp.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code:\ngot %v\nwant %v\n",
				status, http.StatusNotFound)
		}
	}
}