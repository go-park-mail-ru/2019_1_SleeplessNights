package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	requestErrors, isValid, err := router.ValidateRegisterRequest(r)
	if err != nil {
		router.Return500(&w, err)
	}
	if !isValid {
		router.Return400(&w, requestErrors)
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
	salt, err := router.MakeSalt()
	if err != nil {
		router.Return500(&w, err)
		return
	}
	user.Salt = salt
	user.Password = router.MakePasswordHash(r.Form.Get("password"), user.Salt)
	defer func() {
		//Пользователь уже успешно создан, поэтому его в любом случае следует добавить в БД
		//Однако, с ним ещё можно произвести полезную работу, которая может вызвать ошибки
		models.Users[user.Email] = user
		models.UserKeyPairs[user.ID] = user.Email//Пара ключей ID-email, чтобы юзера можно было найти 2-мя способами
	}()

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