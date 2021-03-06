package helpers

import (
	"encoding/json"
	"net/http"
)

const (
	FormParsingErrorMsg         = "Ошибка разбора формы"
	UniqueEmailErrorMsg         = "Пользователь с таким адресом электронной почты уже зарегистрирован"
	MissedUserErrorMsg          = "Неверно введён адрес электронной почты или пароль"
	InvalidEmailErrorMsg        = "Неверно введён адрес электронной почты"
	WrongPassword               = "Неверно введён пароль"
	PasswordsDoNotMatchErrorMsg = "Пароли не совпадают"
	PasswordIsTooSmallErrorMsg  = "Пароль слишком короткий"
	InvalidNicknameErrorMsg     = "Никнейм может состоять только из букв латинского алфавита, цифр и символов '-' и '_'"
	NicknameIsTooSmallErrorMsg  = "Никнейм не может быть короче 4 символов"
	NicknameIsTooLongErrorMsg   = "Никнейм не может быть длиннее 16 символов"
	AvatarExtensionError        = "Файл имеет неподдерживаемый формат"
	AvatarIsMissingError        = "Файл аватара не содержит данных"
	AvatarFileIsTooBig          = "Файл аватара слишком большой (более 2МБайт)"
)

type ErrorSet []string
type errorResponse struct {
	Email     string   `json:"email,omitempty"`
	Password  string   `json:"password,omitempty"`
	Password2 string   `json:"password2,omitempty"`
	Nickname  string   `json:"nickname,omitempty"`
	Avatar    string   `json:"avatar,omitempty"`
	Error     []string `json:"error,omitempty"`
}

func MarshalToJSON(errSet ErrorSet) ([]byte, error) {
	var responseBody errorResponse
	for _, err := range errSet {
		switch err {
		case FormParsingErrorMsg:
			responseBody.Error = append(responseBody.Error, err)
		case UniqueEmailErrorMsg:
			responseBody.Email = err
		case MissedUserErrorMsg:
			responseBody.Error = append(responseBody.Error, err)
		case InvalidEmailErrorMsg:
			responseBody.Email = err
		case WrongPassword:
			responseBody.Password = err
		case PasswordsDoNotMatchErrorMsg:
			responseBody.Password2 = err
		case PasswordIsTooSmallErrorMsg:
			responseBody.Password = err
		case InvalidNicknameErrorMsg:
			responseBody.Nickname = err
		case NicknameIsTooSmallErrorMsg:
			responseBody.Nickname = err
		case NicknameIsTooLongErrorMsg:
			responseBody.Nickname = err
		case AvatarExtensionError:
			responseBody.Avatar = err
		case AvatarIsMissingError:
			responseBody.Avatar = err
		case AvatarFileIsTooBig:
			responseBody.Avatar = err
		default:
			responseBody.Error = append(responseBody.Error, err)
		}
	}
	return json.Marshal(responseBody)
}

func Return500(w *http.ResponseWriter, err error) {
	(*w).WriteHeader(http.StatusInternalServerError)
	data := ErrorSet{err.Error()}
	jsonData, err := MarshalToJSON(data)
	if err != nil {
		logger.Errorf("Error while marshaling json for 500 response")
	}
	_, err = (*w).Write(jsonData)
	if err != nil {
		logger.Errorf("Error while writing request body for 500 response")
	}
}

func Return400(w *http.ResponseWriter, requestErrorMessages ErrorSet) {
	(*w).WriteHeader(http.StatusBadRequest)
	data := requestErrorMessages
	jsonData, err := MarshalToJSON(data)
	if err != nil {
		logger.Errorf("Error while marshaling json for 400 response")
		Return500(w, err)
	}
	_, err = (*w).Write(jsonData)
	if err != nil {
		logger.Errorf("Error while writing request body for 400 response")
		Return500(w, err)
	}
}
