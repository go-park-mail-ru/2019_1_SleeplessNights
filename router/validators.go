package router

import (
	"bytes"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"net/http"
	"regexp"
	"strings"
)

func ValidateRegisterRequest(r *http.Request)(requestErrors ErrorSet, isValid bool, err error) {
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

func ValidateAuthRequest(r *http.Request)(requestErrors ErrorSet, isValid bool, user models.User, err error) {
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

func validateEmail(email string, requestErrors *ErrorSet)(err error) {
	isValid, err := regexp.Match("^[a-z0-9._%+-]+@[a-z0-9-]+.+.[a-z]{2,4}$", []byte(email))
	if !isValid {
		*requestErrors = append(*requestErrors, InvalidEmailErrorMsg)
	}
	return
}

func validatePassword(password string, requestErrors *ErrorSet)(err error) {
	isValid := len(password) >= 8
	if !isValid {
		*requestErrors = append(*requestErrors, PasswordIsTooSmallErrorMsg)
	}
	return nil
}

func validateNickname(nickname string, requestErrors *ErrorSet)(err error) {
	isValid, err := regexp.Match("^[a-z0-9_-]$", []byte(nickname))
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