package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"golang.org/x/net/context"
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

	requestErrors, err := helpers.ValidateAuthRequest(r)
	if requestErrors != nil {
		helpers.Return400(&w, requestErrors)
		return
	}

	sessionToken, err := userManager.MakeToken(context.Background(),
		&services.UserSignature{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		})
	if err != nil {
		switch err.Error() {
		case errors.DataBaseNoDataFound.Error():
			helpers.Return400(&w, helpers.ErrorSet{helpers.MissedUserErrorMsg})
			return
		case errors.AuthWrongPassword.Error():
			helpers.Return400(&w, helpers.ErrorSet{helpers.WrongPassword})
			return
		default:
			helpers.Return500(&w, err)
			return
		}
	}

	sessionCookie := helpers.BuildSessionCookie(sessionToken)
	http.SetCookie(w, &sessionCookie)

	user, err := userManager.CheckToken(context.Background(), sessionToken)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}

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
	cookie, err := r.Cookie(config.GetString("main_ms.pkg.helpers.cookie.name"))
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
