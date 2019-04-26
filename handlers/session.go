package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/auth"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers/helpers"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		formErrorMessages := helpers.ErrorSet{
			helpers.FormParsingErrorMsg,
			err.Error(),
		}
		helpers.Return400(&w, formErrorMessages)
		return
	}

	requestErrors, user, err := helpers.ValidateAuthRequest(r)
	if requestErrors != nil {
		helpers.Return400(&w, requestErrors)
		return
	}

	sessionCookie, err := auth.MakeSession(user)
	if err != nil {
		helpers.Return500(&w, err) //TODO test wrong cookie
		return
	}

	http.SetCookie(w, &sessionCookie)

	data, err := json.Marshal(user)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
}


func AuthDeleteHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie.Value = ""

	http.SetCookie(w, cookie)
	_, err = w.Write([]byte("{}"))
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
}
