package helpers

import (
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	MaxPhotoSize = 2 * 1024 * 1024
)

var avatarTypeWhiteList map[string]struct{}

func ValidateUpdateProfileRequest(r *http.Request) (requestErrors ErrorSet) {
	newNickname := r.Form.Get("nickname")
	validateNickname(newNickname, &requestErrors)

	newEmail := strings.ToLower(r.Form.Get("email"))
	r.Form.Set("email", newEmail)
	validateEmail(newEmail, &requestErrors)

	avatar := r.MultipartForm.File["avatar"][0]
	validateAvatar(avatar, &requestErrors)

	return requestErrors
}

func ValidateRegisterRequest(r *http.Request) (requestErrors ErrorSet) {
	email := strings.ToLower(r.Form.Get("email"))
	r.Form.Set("email", email)
	validateEmail(email, &requestErrors)

	password := r.Form.Get("password")
	validatePassword(password, &requestErrors)

	password2 := r.Form.Get("password2")
	if password != password2 {
		requestErrors = append(requestErrors, PasswordsDoNotMatchErrorMsg)
	}

	nickname := r.Form.Get("nickname")
	validateNickname(nickname, &requestErrors)

	return requestErrors
}

func validateEmail(email string, requestErrors *ErrorSet) {
	isValid := emailReg.Match([]byte(email))
	if !isValid {
		logger.Errorf("Email isn't valid")
		*requestErrors = append(*requestErrors, InvalidEmailErrorMsg)
	}
}

func validatePassword(password string, requestErrors *ErrorSet) {
	isValid := len(password) >= 6
	if !isValid {
		logger.Errorf("Password isn't valid")
		*requestErrors = append(*requestErrors, PasswordIsTooSmallErrorMsg)
	}
}

func validateNickname(nickname string, requestErrors *ErrorSet) {
	isValid := nicknameReg.Match([]byte(nickname))
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
}

func validateAvatar(avatar *multipart.FileHeader, requestErrors *ErrorSet) {

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
