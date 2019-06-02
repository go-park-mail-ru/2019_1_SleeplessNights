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
	profile, err := userManager.GetProfile(context.Background(), user)
	if err != nil {
		logger.Error("Unable to get profile by user:", err)
		helpers.Return500(&w, err)
	}

	data, err := json.Marshal(profile)
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

func ProfileUpdateHandler(user *services.User, w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(helpers.MaxPhotoSize)
	if err != nil {
		formErrorMessages := helpers.ErrorSet{
			helpers.FormParsingErrorMsg,
			err.Error(),
		}
		logger.Errorf("Failed to parse form: %v", err.Error())
		helpers.Return400(&w, formErrorMessages)
		return
	}

	requestErrors := helpers.ValidateUpdateProfileRequest(r)
	if requestErrors != nil {
		logger.Errorf("RequestErrors isn't empty.")
		helpers.Return400(&w, requestErrors)
		return
	}

	user.Nickname = r.MultipartForm.Value["nickname"][0]
	logger.Info("nickname ", user.Nickname)
	newAvatar := r.MultipartForm.File["avatar"][0]
	logger.Info("newAvatar ", newAvatar)
	avatarName := uuid.NewV4().String() + filepath.Ext(newAvatar.Filename)
	user.AvatarPath = avatarName
	_, err = userManager.UpdateProfile(context.Background(), user)
	if err != nil {
		logger.Errorf("Failed to update profile: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}

	avatarFile, err := newAvatar.Open()
	if err != nil {
		logger.Errorf("Failed to open file: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}

	defer func() {
		err := avatarFile.Close()
		if err != nil {
			logger.Errorf("Failed to close file: %v", err.Error())
			helpers.Return500(&w, err)
			return
		}
	}()

	avatarBytes, err := ioutil.ReadAll(avatarFile)
	logger.Info("Avatar bytes :", avatarBytes)
	if err != nil {
		logger.Errorf("Failed to read all file: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}
	logger.Info(" New Image is saved to \"" + os.Getenv("BASEPATH") + AvatarPrefix + avatarName + "\" ")
	file, err := os.Create(os.Getenv("BASEPATH") + AvatarPrefix + avatarName)
	if err != nil {
		logger.Errorf("Failed to create file: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}

	defer func() {
		err := file.Close()
		if err != nil {
			logger.Errorf("Failed to close file: %v", err.Error())
			helpers.Return500(&w, err)
			return
		}
	}()

	_, err = file.Write(avatarBytes)
	if err != nil {
		logger.Errorf("Failed to write response: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}

	_, err = w.Write([]byte(`{"avatarPath": "` + avatarName + `"}`))
	if err != nil {
		logger.Errorf("Failed to write response: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}
}
