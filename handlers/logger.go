package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	LoggerPath = "logger/logBase.log/"
)

func LoggerHandler(w http.ResponseWriter, r *http.Request) {
	router.SetBasicHeaders(&w)

	_, err := os.Stat(LoggerPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	logger, err := ioutil.ReadFile(LoggerPath)
	if err != nil {
		router.Return500(&w, err)
		return
	}

	w.Header().Set("Content-type", http.DetectContentType(logger))
	_, err = w.Write(logger)
	if err != nil {
		router.Return500(&w, err)
		return
	}
}
