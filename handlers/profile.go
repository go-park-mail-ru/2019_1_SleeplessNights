package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"io/ioutil"
	"net/http"
	"os"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte("{}"))
		if err != nil {
			helpers.Return500(&w, err)
			return
		}
		return
	}
	user, err := helpers.Authorize(sessionCookie.Value)
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte("{}"))
		if err != nil {
			helpers.Return500(&w, err)
			return
		}
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

func ProfileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(helpers.MaxPhotoSize)
	if err != nil {
		formErrorMessages := helpers.ErrorSet{
			helpers.FormParsingErrorMsg,
			err.Error(),
		}
		helpers.Return400(&w, formErrorMessages)
		return
	}

	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	/*requestErrors, isValid, err := helpers.ValidateUpdateProfileRequest(r, user) //TODO WRITE VALIDATOR
	if err != nil {
		helpers.Return500(&w, err)
	}
	if !isValid {
		helpers.Return400(&w, requestErrors)
		return
	}*/

	newAvatar:= r.MultipartForm.File["avatar"][0]
	avatarFile, err := newAvatar.Open()
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	defer func() {
		err := avatarFile.Close()
		if err != nil {
			helpers.Return500(&w, err)
			return
		}
	}()

	avatarBytes, err := ioutil.ReadAll(avatarFile)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	//TODO rename avatar name in server
	newAvatarName := r.MultipartForm.File["avatar"][0].Filename
	file, err := os.Create(avatarPrefix + newAvatarName)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	defer func(){
		err := file.Close()
		if err != nil {
			helpers.Return500(&w, err)
			return
		}
	}()

	_, err = file.Write(avatarBytes)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}

	user, err := helpers.Authorize(sessionCookie.Value)
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user.Nickname = r.MultipartForm.Value["nickname"][0]
	user.AvatarPath = newAvatarName
	models.Users[user.Email] = user
	_, err = w.Write([]byte(`{"avatar_path": "`+newAvatarName+`"}`))
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
}
