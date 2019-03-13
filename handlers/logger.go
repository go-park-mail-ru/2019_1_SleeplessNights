package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"io/ioutil"
	"net/http"
	"os"
)

func LoggerHandler(w http.ResponseWriter, r *http.Request) {
	router.SetBasicHeaders(&w)

	_, err := os.Stat(logger.LoggerPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_logger, err := ioutil.ReadFile(logger.LoggerPath)
	if err != nil {
		router.Return500(&w, err)
		return
	}

	w.Header().Set("Content-type", "text/plain")
	_, err = w.Write(_logger)
	if err != nil {
		router.Return500(&w, err)
		return
	}
}
