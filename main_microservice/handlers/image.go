package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
)

var AvatarPrefix = config.GetString("main_ms.pkg.handlers.avatar_prefix")

func ImgHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pathToFile, found := vars["path"]
	if !found {
		logger.Errorf("Didn't find `path`.")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	path := os.Getenv("BASEPATH") + AvatarPrefix + pathToFile

	_, err := os.Stat(path)
	if err != nil {
		logger.Errorf("Failed to stat path: %v", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	avatar, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Errorf("Failed to read file: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}
	w.Header().Set("Content-type", http.DetectContentType(avatar))
	_, err = w.Write(avatar)
	if err != nil {
		logger.Errorf("Failed to write response: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}
}
