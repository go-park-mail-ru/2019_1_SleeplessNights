package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"net/http"
)

func OptionsHandler(w http.ResponseWriter, r *http.Request) {
	router.SetBasicHeaders(&w)
	w.WriteHeader(http.StatusNoContent)
}
