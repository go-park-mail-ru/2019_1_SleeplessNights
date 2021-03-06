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
		logger.Errorf("Failed to parse form: %v", err.Error())
		helpers.Return400(&w, formErrorMessages)
		return
	}

	sessionToken, err := userManager.MakeToken(context.Background(),
		&services.UserSignature{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		})
	if err != nil {
		logger.Errorf("Failed to make token: %v", err.Error())
		matchedUV:= errors.DataBaseUniqueViolationReg.Match([]byte(err.Error()))
		matchedNDF:= errors.DataBaseNoDataFoundReg.Match([]byte(err.Error()))
		if matchedUV {
			helpers.Return400(&w, helpers.ErrorSet{helpers.MissedUserErrorMsg})
			return
		} else if matchedNDF {
			helpers.Return400(&w, helpers.ErrorSet{helpers.MissedUserErrorMsg})
			return
		} else {
			helpers.Return500(&w, err)
			return
		}
	}

	sessionCookie := helpers.BuildSessionCookie(sessionToken)
	http.SetCookie(w, &sessionCookie)

	user, err := userManager.CheckToken(context.Background(), sessionToken)
	if err != nil {
		logger.Errorf("Failed to check token: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		logger.Errorf("Failed to marshal user: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		logger.Errorf("Failed to write response: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}
}

func AuthDeleteHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(config.GetString("main_ms.pkg.helpers.cookie.name"))
	if err != nil {
		logger.Errorf("Failed to get cookie: %v", err.Error())
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie.Value = ""

	http.SetCookie(w, cookie)
	_, err = w.Write([]byte("{}"))
	if err != nil {
		logger.Errorf("Failed to write response: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}
}
