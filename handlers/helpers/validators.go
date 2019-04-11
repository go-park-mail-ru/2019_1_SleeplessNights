package helpers

import (
	"bytes"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

const (
	MaxPhotoSize = 2 * 1024 * 1024
)

var avatarTypeWhiteList map[string]struct{}

func ValidateUpdateProfileRequest(r *http.Request) (requestErrors ErrorSet, isValid bool, err error) {
	newNickname := r.Form.Get("nickname")

	err = validateNickname(newNickname, &requestErrors)
	if err != nil {
		logger.Error.Println("Failed to update profile:", err)
		return
	}

	newEmail := strings.ToLower(r.Form.Get("email"))
	r.Form.Set("email", newEmail)

	err = validateEmail(newEmail, &requestErrors)
	if err != nil {
		logger.Error.Println("Failed to update profile:", err)
		return
	}

	avatar := r.MultipartForm.File["avatar"][0]

	err = validateAvatar(avatar, &requestErrors)
	if err != nil {
		logger.Error.Println("Failed to update profile:", err)
		return
	}
	return requestErrors, len(requestErrors) == 0, nil
}

func ValidateRegisterRequest(r *http.Request) (requestErrors ErrorSet, err error) {
	email := strings.ToLower(r.Form.Get("email"))
	r.Form.Set("email", email)
	err = validateEmail(email, &requestErrors)
	if err != nil {
		return
	}

	password := r.Form.Get("password")
	err = validatePassword(password, &requestErrors)
	if err != nil {
		return
	}

	password2 := r.Form.Get("password2")
	if password != password2 {
		requestErrors = append(requestErrors, PasswordsDoNotMatchErrorMsg)
	}

	nickname := r.Form.Get("nickname")
	err = validateNickname(nickname, &requestErrors)
	if err != nil {
		return
	}
	return requestErrors, nil
}

func ValidateAuthRequest(r *http.Request) (requestErrors ErrorSet, user models.User, err error) {
	email := strings.ToLower(r.Form.Get("email"))
	r.Form.Set("email", email)
	err = validateEmail(email, &requestErrors)
	if err != nil {
		return
	}

	password := r.Form.Get("password")
	err = validatePassword(password, &requestErrors)
	if err != nil {
		return
	}

	user, err = database.GetInstance().GetUserViaEmail(email)
	if err != nil {
		if err.Error() == SQLNoRows{
			requestErrors = append(requestErrors, MissedUserErrorMsg)
			return
		} else {
			return
		}
	}

	hashedPassword := MakePasswordHash(password, user.Salt)
	if bytes.Compare(hashedPassword, user.Password) != 0 {
		requestErrors = append(requestErrors, WrongPassword)
	}

	return requestErrors, user, nil
}

func validateEmail(email string, requestErrors *ErrorSet) (err error) {
	isValid, err := regexp.Match("^[a-z0-9._%+-]+@[a-z0-9-]+.+.[a-z]{2,4}$", []byte(email))
	if !isValid {
		*requestErrors = append(*requestErrors, InvalidEmailErrorMsg)
	}
	return
}

func validatePassword(password string, requestErrors *ErrorSet) (err error) {
	isValid := len(password) >= 8
	if !isValid {
		*requestErrors = append(*requestErrors, PasswordIsTooSmallErrorMsg)
	}
	return nil
}

func validateNickname(nickname string, requestErrors *ErrorSet) (err error) {
	isValid, err := regexp.Match("^[A-Za-z0-9_-]", []byte(nickname))
	if err != nil {
		return
	}
	if !isValid {
		*requestErrors = append(*requestErrors, InvalidNicknameErrorMsg)
	}
	if len(nickname) < 3 {
		*requestErrors = append(*requestErrors, NicknameIsTooSmallErrorMsg)
	}
	if len(nickname) > 16 {
		*requestErrors = append(*requestErrors, NicknameIsTooLongErrorMsg)
	}
	return
}

func validateAvatar(avatar *multipart.FileHeader, requestErrors *ErrorSet) (err error) {

	if avatar.Size == 0 {
		*requestErrors = append(*requestErrors, AvatarIsMissingError)
	}
	if avatar.Size > MaxPhotoSize {
		*requestErrors = append(*requestErrors, AvatarFileIsTooBig)
	}
	contentType := avatar.Header.Get("content-type")
	if _, found := avatarTypeWhiteList[contentType]; !found {
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
