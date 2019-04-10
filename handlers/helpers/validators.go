package helpers

import (
	"bytes"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/lib/pq"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

const (
	MaxPhotoSize = 2 * 1024 * 1024
)

var avatarTypeWhiteList map[string]struct{}

func ValidateUpdateProfileRequest(r *http.Request, user models.User) (requestErrors ErrorSet, isValid bool, err error) {
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

	existingUser, userFound, err := database.GetInstance().GetUserViaEmail(newEmail)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Code.Class())
		logger.Error.Print(_err.Error())
		return
	}
	if  userFound && user.ID != existingUser.ID {
		logger.Error.Println("Failed to update profile:", UniqueEmailErrorMsg)
		requestErrors = append(requestErrors, UniqueEmailErrorMsg)
	}

	avatar := r.MultipartForm.File["avatar"][0]

	err = validateAvatar(avatar, &requestErrors)
	if err != nil {
		logger.Error.Println("Failed to update profile:", err)
		return
	}
	return requestErrors, len(requestErrors) == 0, nil
}

func ValidateRegisterRequest(r *http.Request) (requestErrors ErrorSet, isValid bool, err error) {
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

	_, userExist, err := database.GetInstance().GetUserViaEmail(r.Form.Get("email"))
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Code.Class())
		logger.Error.Print(_err.Error())
		return
	}
	if userExist {
		requestErrors = append(requestErrors, UniqueEmailErrorMsg)
		return
	}

	return requestErrors, len(requestErrors) == 0, nil
}

func ValidateAuthRequest(r *http.Request) (requestErrors ErrorSet, isValid bool, user models.User, err error) {
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

	user, found, err := database.GetInstance().GetUserViaEmail(email)
	if _err, ok := err.(*pq.Error); ok {
		logger.Error.Print(_err.Code.Class())
		logger.Error.Print(_err.Error())
		return
	}
	if !found {
		requestErrors = append(requestErrors, MissedUserErrorMsg)
		return requestErrors, false, user, nil
	}

	hashedPassword := MakePasswordHash(password, user.Salt)
	if bytes.Compare(hashedPassword, user.Password) != 0 {
		requestErrors = append(requestErrors, WrongPassword)
	}

	return requestErrors, len(requestErrors) == 0, user, nil
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
