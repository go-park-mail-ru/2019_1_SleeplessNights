package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers/helpers"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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


	user, err := helpers.Authorize(sessionCookie.Value)
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	requestErrors, err := helpers.ValidateUpdateProfileRequest(r)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	if requestErrors != nil {
		helpers.Return400(&w, requestErrors)
		return
	}

	user.Nickname = r.MultipartForm.Value["nickname"][0]
	newEmail := r.MultipartForm.Value["email"][0]
	userID := user.ID
	user.Email = newEmail

	newAvatar := r.MultipartForm.File["avatar"][0]
	avatarName := uuid.NewV4().String() + filepath.Ext(newAvatar.Filename)
	user.AvatarPath = avatarName


	err = database.GetInstance().UpdateUser(user, userID)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}

	cookie, err := helpers.MakeSession(user)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	http.SetCookie(w, &cookie)

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

	file, err := os.Create(os.Getenv("BASEPATH") + AvatarPrefix + avatarName)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}

	defer func() {
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

	_, err = w.Write([]byte(`{"avatar_path": "` + avatarName + `"}`))
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
}
