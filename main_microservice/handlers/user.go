package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/auth_microservice"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/models"
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

	requestErrors, err := helpers.ValidateRegisterRequest(r)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	if requestErrors != nil {
		helpers.Return400(&w, requestErrors)
		return
	}

	user := models.User{
		Email:      r.Form.Get("email"),
		Won:        0,
		Lost:       0,
		PlayTime:   0,
		Nickname:   r.Form.Get("nickname"),
		AvatarPath: "default_avatar.jpg",
	}
	salt, err := helpers.MakeSalt()
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	user.Salt = salt
	user.Password = helpers.MakePasswordHash(r.Form.Get("password"), user.Salt)

	err = database.GetInstance().AddUser(user)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}

	user, err = database.GetInstance().GetUserViaEmail(user.Email)//id присваивается только при добавлении в базу
	if err != nil {
		helpers.Return500(&w, err)
		return
	}

	sessionCookie, err := auth_microservice.MakeSession(user)
	if err != nil {
		helpers.Return500(&w, err)
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
