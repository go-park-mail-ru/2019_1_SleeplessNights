package helpers

import (
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

const (
	MaxPhotoSize = 2 * 1024 * 1024
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("MM_Validator")
	logger.SetLogLevel(logrus.ErrorLevel)
}

var avatarTypeWhiteList map[string]struct{}

func ValidateUpdateProfileRequest(r *http.Request) (requestErrors ErrorSet, err error) {
	newNickname := r.Form.Get("nickname")

	err = validateNickname(newNickname, &requestErrors)
	if err != nil {
		logger.Error("Failed to update profile:", err)
		return
	}

	newEmail := strings.ToLower(r.Form.Get("email"))
	r.Form.Set("email", newEmail)

	err = validateEmail(newEmail, &requestErrors)
	if err != nil {
		logger.Error("Failed to update profile:", err)
		return
	}

	avatar := r.MultipartForm.File["avatar"][0]

	err = validateAvatar(avatar, &requestErrors)
	if err != nil {
		logger.Error("Failed to update profile:", err)
		return
	}
	return requestErrors, nil
}

func ValidateRegisterRequest(r *http.Request) (requestErrors ErrorSet, err error) {
	email := strings.ToLower(r.Form.Get("email"))
	r.Form.Set("email", email)
	err = validateEmail(email, &requestErrors)
	if err != nil {
		logger.Errorf("Failed to validate email: %v", err.Error())
		return
	}

	password := r.Form.Get("password")
	err = validatePassword(password, &requestErrors)
	if err != nil {
		logger.Errorf("Failed to get form `password`: %v", err.Error())
		return
	}

	password2 := r.Form.Get("password2")
	if password != password2 {
		logger.Errorf("Failed to get form `password2`: %v", err.Error())
		requestErrors = append(requestErrors, PasswordsDoNotMatchErrorMsg)
	}

	nickname := r.Form.Get("nickname")
	err = validateNickname(nickname, &requestErrors)
	if err != nil {
		logger.Errorf("Failed to get form `nickname`: %v", err.Error())
		return
	}
	return requestErrors, nil
}

func validateEmail(email string, requestErrors *ErrorSet) (err error) {
	isValid, err := regexp.Match("^[a-z0-9._%+-]+@[a-z0-9-]+.+.[a-z]{2,4}$", []byte(email))
	if !isValid {
		logger.Errorf("Email isn't valid")
		*requestErrors = append(*requestErrors, InvalidEmailErrorMsg)
	}
	return
}

func validatePassword(password string, requestErrors *ErrorSet) (err error) {
	isValid := len(password) >= 6
	if !isValid {
		logger.Errorf("Password isn't valid")
		*requestErrors = append(*requestErrors, PasswordIsTooSmallErrorMsg)
	}
	return nil
}

func validateNickname(nickname string, requestErrors *ErrorSet) (err error) {
	isValid, err := regexp.Match("^[a-zA-Z0-9-_]*$", []byte(nickname))
	if err != nil {
		logger.Errorf("Failed to match: %v", err.Error())
		return
	}
	if !isValid {
		logger.Errorf("Nickname isn't valid")
		*requestErrors = append(*requestErrors, InvalidNicknameErrorMsg)
	}
	if len(nickname) <= 3 {
		logger.Errorf("Nickname is short")
		*requestErrors = append(*requestErrors, NicknameIsTooSmallErrorMsg)
	}
	if len(nickname) >= 17 {
		logger.Errorf("Nickname is long")
		*requestErrors = append(*requestErrors, NicknameIsTooLongErrorMsg)
	}
	return
}

func validateAvatar(avatar *multipart.FileHeader, requestErrors *ErrorSet) (err error) {

	if avatar.Size == 0 {
		logger.Errorf("Avatar is empty")
		*requestErrors = append(*requestErrors, AvatarIsMissingError)
	}
	if avatar.Size > MaxPhotoSize {
		logger.Errorf("Avatar is big")
		*requestErrors = append(*requestErrors, AvatarFileIsTooBig)
	}
	contentType := avatar.Header.Get("content-type")
	if _, found := avatarTypeWhiteList[contentType]; !found {
		logger.Errorf("Avatar didn't find")
		*requestErrors = append(*requestErrors, AvatarExtensionError)
	}

	return nil
}
func init() {
	avatarTypeWhiteList = make(map[string]struct{})
	avatarTypeWhiteList["image/gif"] = struct{}{}
	avatarTypeWhiteList["image/png"] = struct{}{}
	avatarTypeWhiteList["image/jpeg"] = struct{}{}
	avatarTypeWhiteList["image/bmp"] = struct{}{}
	avatarTypeWhiteList["image/tiff"] = struct{}{}
	avatarTypeWhiteList["image/pjpeg"] = struct{}{}
}
