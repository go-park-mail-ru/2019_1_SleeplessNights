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
		helpers.Return400(&w, formErrorMessages)
		return
	}
	logger.Debug("ParseForm_OK")
	requestErrors, err := helpers.ValidateRegisterRequest(r)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	if requestErrors != nil {
		helpers.Return400(&w, requestErrors)
		return
	}
	logger.Debug("ValidateRegisterRequest_OK")

	user, err := userManager.CreateUser(context.Background(),
		&services.NewUserData{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
			Nickname: r.Form.Get("nickname"),
		})
	if err != nil{
		switch err.Error() {
		case errors.DataBaseUniqueViolation.Error():
			helpers.Return400(&w, helpers.ErrorSet{helpers.UniqueEmailErrorMsg})
			return
		default:
			helpers.Return500(&w, err)
			return
		}
	}
	logger.Debug("CreateUser_OK")

	sessionToken, err := userManager.MakeToken(context.Background(),
		&services.UserSignature{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		})
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	logger.Debug("MakeToken_OK")


	sessionCookie := helpers.BuildSessionCookie(sessionToken)
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
	logger.Debug("BuildSessionCookie_OK")
}
