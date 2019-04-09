package handlers_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/faker"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
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

	faker.CreateFakeData(handlers.UserCounter)

	for _, user := range database.GetInstance().GetUsers() {
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

	form := url.Values{}
	form.Add("WRONG_mail", "test@test.com")
	form.Add("WRONG_password", "asdasdasdsadasd")

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
