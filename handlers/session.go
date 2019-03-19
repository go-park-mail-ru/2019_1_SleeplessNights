package handlers

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/handlers/helpers"
	"net/http"
)

func AuthHandler(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		formErrorMessages := helpers.ErrorSet{
			helpers.FormParsingErrorMsg,
			err.Error(),
		}
		helpers.Return400(&w, formErrorMessages)
		return
	}

	requestErrors, isValid, user, err := helpers.ValidateAuthRequest(r)
	if err != nil {
		helpers.Return500(&w, err)
	}
	if !isValid {
		helpers.Return400(&w, requestErrors)
		return
	}

	sessionCookie, err := helpers.MakeSession(user)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	http.SetCookie(w, &sessionCookie)
	_, err = w.Write([]byte("{}"))
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
}