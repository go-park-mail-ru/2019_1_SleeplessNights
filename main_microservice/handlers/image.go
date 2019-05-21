package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	AvatarPrefix = "/main_microservice/static/img/"
)

func ImgHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pathToFile, found := vars["path"]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	path := os.Getenv("BASEPATH") + AvatarPrefix + pathToFile

	_, err := os.Stat(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	avatar, err := ioutil.ReadFile(path)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	content_type := http.DetectContentType(avatar)
	if content_type == "text/plain; charset=utf-8" {
		w.Header().Set("Content-type", "image/svg+xml")
	} else {
		w.Header().Set("Content-type", content_type)
	}

	_, err = w.Write(avatar)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
}
