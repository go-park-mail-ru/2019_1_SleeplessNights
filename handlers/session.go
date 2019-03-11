package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request){
	router.SetBasicHeaders(&w)
	err := r.ParseForm()
	if err != nil {
		formErrorMessages := router.ErrorSet{
			router.FormParsingErrorMsg,
			err.Error(),
		}
		router.Return400(&w, formErrorMessages)
		return
	}

	requestErrors, isValid, user, err := router.ValidateAuthRequest(r)
	if err != nil {
		router.Return500(&w, err)
	}
	if !isValid {
		router.Return400(&w, requestErrors)
		return
	}

	sessionCookie, err := router.MakeSession(user)
	if err != nil {
		router.Return500(&w, err)
		return
	}
	http.SetCookie(w, &sessionCookie)
	_, err = w.Write([]byte("{}"))
	if err != nil {
		router.Return500(&w, err)
		return
	}
}