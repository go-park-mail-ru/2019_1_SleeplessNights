package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"golang.org/x/net/context"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	requestErrors, err := helpers.ValidateRegisterRequest(r)
	if err != nil {
		logger.Errorf("Failed to validate request: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}
	if requestErrors != nil {
		logger.Errorf("RequestErrors isn't empty.")
		helpers.Return400(&w, requestErrors)
		return
	}

	user, err := userManager.CreateUser(context.Background(),
		&services.NewUserData{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
			Nickname: r.Form.Get("nickname"),
		})
	if err != nil {
		logger.Errorf("Failed to create user: %v", err.Error())
		matchedUV := errors.DataBaseUniqueViolationReg.Match([]byte(err.Error()))
		if matchedUV {
			helpers.Return400(&w, helpers.ErrorSet{helpers.UniqueEmailErrorMsg})
			return
		} else {
			helpers.Return500(&w, err)
			return
		}
	}

	sessionToken, err := userManager.MakeToken(context.Background(),
		&services.UserSignature{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		})
	if err != nil {
		logger.Errorf("Failed to make token: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}

	sessionCookie := helpers.BuildSessionCookie(sessionToken)
	http.SetCookie(w, &sessionCookie)

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
