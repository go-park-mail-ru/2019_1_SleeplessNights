package helpers

import (
	"bytes"
	"errors"
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

func ValidateUpdateProfileRequest(r *http.Request, user models.User) (requestErrors ErrorSet, isValid bool, err error) {

	newNickname := strings.ToLower(r.MultipartForm.Value["nickname"][0])
	err = validateNickname(newNickname, &requestErrors)

	if err != nil {
		logger.Error.Println("Failed to update profile:", err)
		return
	}

	newEmail := strings.ToLower(r.MultipartForm.Value["email"][0])
	err = validateEmail(newEmail, &requestErrors)

	if err != nil {
		logger.Error.Println("Failed to update profile:", err)
		return
	}

	if existingUser, userFound := models.Users[newEmail]; userFound && user.ID != existingUser.ID {
		logger.Error.Println("Failed to update profile:", UniqueEmailErrorMsg)
		requestErrors = append(requestErrors, UniqueEmailErrorMsg)
		return
	}

	avatar := r.MultipartForm.File["avatar"][0]

	err = validateAvatar(avatar, &requestErrors)

	if err != nil {
		logger.Error.Println("Failed to update profile:", err)
		return
	}

	isValid = true
	return
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

	_, userExist := models.Users[r.Form.Get("email")]
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

	user, found := models.Users[email]
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
		return errors.New("файл нового аватара пустой")
	}
	//File is bigger than 10 MBytes
	if avatar.Size > 10e6 {
		*requestErrors = append(*requestErrors, AvatarFileIsTooBig)
		return errors.New("файл аватара слишком большой")
	}

	if contentType := avatar.Header.Get("content-type"); contentType != "image/jpeg" {
		*requestErrors = append(*requestErrors, AvatarExtensionError)
		return errors.New("неизвестный формат файла")
	}

	return
}
