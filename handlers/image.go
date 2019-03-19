package handlers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers/helpers"
	"github.com/gorilla/mux"

	"io/ioutil"
	"net/http"
	"os"
)

const (
	AvatarPrefix = "../static/img/"
)

func ImgHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pathToFile, found := vars["path"]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	path := AvatarPrefix + pathToFile
	_, err := os.Stat(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Printf("!!!!\n")
	avatar, err := ioutil.ReadFile(path)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	w.Header().Set("Content-type", http.DetectContentType(avatar))
	_, err = w.Write(avatar)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
}
