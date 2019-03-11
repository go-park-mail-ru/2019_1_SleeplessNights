package router

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	FormParsingErrorMsg         = "Ошибка разбора формы"
	UniqueEmailErrorMsg         = "Пользователь с таким адресом электронной почты уже зарегистрирован"
	MissedUserErrorMsg          = "Пользователь с таким адресом электронной почты не зарегистрирован"
	InvalidEmailErrorMsg        = "Неверно введён адрес электронной почты"
	WrongPassword               = "Неверно введён пароль"
	PasswordsDoNotMatchErrorMsg = "Пароли не совпадают"
	PasswordIsTooSmallErrorMsg  = "Пароль слишком короткий"
	InvalidNicknameErrorMsg     = "Никнейм может состоять только из букв латинского алфавита и символов '-' и '_'"
	NicknameIsTooSmallErrorMsg  = "Никнейм не может быть короче 3 символов"
	NicknameIsTooLongErrorMsg   = "Никнейм не может быть длиннее 16 символов"
)

const (
	NoTokenOwner      ="error: There are no token's owner in database"
)

type ErrorSet []string
type errorResponse struct {
	Email     string   `json:"email, omitempty"`
	Password  string   `json:"password, omitempty"`
	Password2 string   `json:"password2, omitempty"`
	Nickname  string   `json:"nickname, omitempty"`
	Error     []string `json:"error, omitempty"`
}

func MarshalToJSON(errSet ErrorSet)([]byte, error) {
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
		log.Println("Error while marshaling json for 500 response")
	}
	_, err = (*w).Write(jsonData)
	if err != nil {
		log.Println("Error while writing request body for 500 response")
	}
}

func Return400(w *http.ResponseWriter, requestErrorMessages ErrorSet) {
	(*w).WriteHeader(http.StatusBadRequest)
	data := requestErrorMessages
	jsonData, err := MarshalToJSON(data)
	if err != nil {
		log.Println("Error while marshaling json for 400 response")
		Return500(w, err)
	}
	_, err = (*w).Write(jsonData)
	if err != nil {
		log.Println("Error while writing request body for 400 response")
		Return500(w, err)
	}
}