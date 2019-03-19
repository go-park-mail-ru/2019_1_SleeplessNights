package handlers

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/handlers/helpers"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/models"
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

	requestErrors, isValid, err := helpers.ValidateRegisterRequest(r)
	if err != nil {
		helpers.Return500(&w, err)
	}
	if !isValid {
		helpers.Return400(&w, requestErrors)
		return
	}

	user := models.User{
		ID:        models.MakeID(),
		Email:     r.Form.Get("email"),
		Won:       0,
		Lost:      0,
		PlayTime:  0,
		Nickname: r.Form.Get("nickname"),
		AvatarPath: "default_avatar.jpg",
	}
	salt, err := helpers.MakeSalt()
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	user.Salt = salt
	user.Password = helpers.MakePasswordHash(r.Form.Get("password"), user.Salt)
	defer func() {
		//Пользователь уже успешно создан, поэтому его в любом случае следует добавить в БД
		//Однако, с ним ещё можно произвести полезную работу, которая может вызвать ошибки
		models.Users[user.Email] = user
		models.UserKeyPairs[user.ID] = user.Email//Пара ключей ID-email, чтобы юзера можно было найти 2-мя способами
	}()

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