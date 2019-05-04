package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func ProfileHandler(user *services.User, w http.ResponseWriter, r *http.Request) {
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

func ProfileUpdateHandler(user *services.User, w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(helpers.MaxPhotoSize)
	if err != nil {
		formErrorMessages := helpers.ErrorSet{
			helpers.FormParsingErrorMsg,
			err.Error(),
		}
		helpers.Return400(&w, formErrorMessages)
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
	newAvatar := r.MultipartForm.File["avatar"][0]
	avatarName := uuid.NewV4().String() + filepath.Ext(newAvatar.Filename)
	user.AvatarPath = avatarName

	_, err = userManager.UpdateProfile(context.Background(), &user)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}

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
