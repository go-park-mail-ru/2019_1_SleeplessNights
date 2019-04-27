package handlers_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/faker"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestCaseAuth struct {
	number   int
	email    string
	password string
	error    string
}

func TestAuthHandlerSuccessfulWithCreateFakeData(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	faker.CreateFakeData(handlers.UserCounter)

	users, err := database.GetInstance().GetUsers()
	if err != nil {
		t.Error(err.Error())
	}
	for _, user := range users {

		email := user.Email
		password := faker.FakeUserPassword

		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)

		req := httptest.NewRequest(http.MethodPost, handlers.ApiAuth, nil)
		req.PostForm = form

		resp := httptest.NewRecorder()

		http.HandlerFunc(handlers.AuthHandler).ServeHTTP(resp, req)

		if status := resp.Code; status == http.StatusInternalServerError {
			t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't make cookie or can't check validate",
				status)
		} else {
			if status := resp.Code; status != http.StatusOK {
				t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v;\n error: %s\n",
					status, http.StatusOK, resp.Body.String())
			}

			expected := `{}`
			if resp.Body.String() != expected {
				t.Errorf("\nhandler returned unexpected body:\ngot %v\nwant %v\nemail: %s\nnickname: %s\npassword: %s\n ",
					resp.Body.String(), expected, user.Email, user.Nickname, string(user.Password))
			}
		}
	}
}

func TestAuthHandlerUnsuccessfulWrongFormsAndNotRegister(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	cases := []TestCaseAuth{
		TestCaseAuth{
			number:   1,
			email:    "test@1.com", //TODO with only numbers after @
			password: "asdasdasdsadasdQ",
			error:    `{"email":"","password":"","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		},
		TestCaseAuth{
			number:   2,
			email:    "test@@test.com",
			password: "asdasdasdsadasdQ",
			error:    `{"email":"Неверно введён адрес электронной почты","password":"","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		},
		TestCaseAuth{
			number:   3,
			email:    "te&st@test.com",
			password: "asdasdasdsadasdQ",
			error:    `{"email":"Неверно введён адрес электронной почты","password":"","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		},
		TestCaseAuth{
			number:   4,
			email:    "test@test.com",
			password: "asdsdQ",
			error:    `{"email":"","password":"Пароль слишком короткий","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		},
		TestCaseAuth{
			number:   5,
			email:    "teesttest.com",
			password: "asdasdasdsadasdQ",
			error:    `{"email":"Неверно введён адрес электронной почты","password":"","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		},
		TestCaseAuth{
			number:   6,
			email:    "_____@test.com",
			password: "asdasdasdsadasd",
			error:    `{"email":"","password":"","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		},
		TestCaseAuth{
			number:   7,
			email:    "teest@test.com",
			password: "",
			error:    `{"email":"","password":"Пароль слишком короткий","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		},
		TestCaseAuth{
			number:   8,
			email:    "",
			password: "asdasdasdsadasd",
			error:    `{"email":"Неверно введён адрес электронной почты","password":"","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		},
		TestCaseAuth{
			number:   9,
			email:    "test@testcom", //TODO without point
			password: "asdasdasdsadasd",
			error:    `{"email":"","password":"","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		},
		TestCaseAuth{
			number:   10,
			email:    "test@test.com",
			password: "asdasdasdsadasd",
			error:    `{"email":"","password":"","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		},
	}

	for _, item := range cases {
		email := item.email
		password := item.password

		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)

		req := httptest.NewRequest(http.MethodPost, handlers.ApiAuth, nil)
		req.PostForm = form

		resp := httptest.NewRecorder()

		http.HandlerFunc(handlers.AuthHandler).ServeHTTP(resp, req)

		if status := resp.Code; status == http.StatusInternalServerError {
			t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't make cookie or can't check validate",
				status)
		} else {
			if status := resp.Code; status != http.StatusBadRequest {
				t.Errorf("\ntest %d.\nhandler returned wrong status code:\ngot %v\nwant %v;\n",
					item.number, status, http.StatusBadRequest)
			}

			if response := resp.Body.String(); response != item.error {
				t.Errorf("\ntest %d.\nhandler returned wrong error response:\ngot %v\nwant %v;\n",
					item.number, response, item.error)
			}
		}
	}
}

func TestAuthHandlerUnsuccessfulWrongParseForm(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	form := url.Values{}
	form.Add("WRONG_mail", "test@test.com")
	form.Add("WRONG_password", "asasdsadasd")

	req := httptest.NewRequest(http.MethodPost, handlers.ApiAuth, nil)
	req.PostForm = form

	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.AuthHandler).ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't make cookie or can't check validate",
			status)
	} else {
		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v;\n",
				status, http.StatusBadRequest)
		}

		//TODO change this expected
		expected := `{"email":"Неверно введён адрес электронной почты","password":"Пароль слишком короткий","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`
		if response := resp.Body.String(); response != expected {
			t.Errorf("\nhandler returned wrong error response:\ngot %v\nwant %v;\n",
				response, expected)
		}
	}
}

func TestAuthDeleteHandlerSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	user := models.User{
		ID:         1,
		Email:      "first@mail.com",
		Nickname:   "first",
		AvatarPath: "none",
	}
	err = database.GetInstance().AddUser(user)
	if err != nil {
		t.Error(err.Error())
	}

	cookie, err := helpers.BuildSessionCookie(user.ID)
	if err != nil {
		t.Errorf("MakeSession returned error: %s\n", err.Error())
		return
	}

	req := httptest.NewRequest(http.MethodDelete, handlers.ApiProfile, nil)
	req.AddCookie(&cookie)

	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.AuthDeleteHandler).ServeHTTP(resp, req)

	expected := "session_token="
	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce",
			status)
	} else {
		if status := resp.Code; status != http.StatusOK {
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v;\n",
				status, http.StatusOK)
		}

		if cookie := resp.Header().Get("Set-Cookie"); cookie != expected {
			t.Errorf("\nhandler returned wrong cookie:\ngot %v\nwant %v;\n",
				cookie, expected)
		}
	}
}

func TestAuthDeleteHandlerUnsuccessful(t *testing.T) {

	req := httptest.NewRequest(http.MethodDelete, handlers.ApiProfile, nil)

	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.AuthDeleteHandler).ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce",
			status)
	} else {
		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v;\n",
				status, http.StatusBadRequest)
		}
	}
}
