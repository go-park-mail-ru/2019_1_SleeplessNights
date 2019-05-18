package handlers_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestCaseAuth struct {
	number    int
	email     string
	nickname  string
	password1 string
	password2 string
	error     string
}

func TestAuthHandlerSuccessfulWithCreateFakeData(t *testing.T) {

	var userManager services.UserMSClient
	var err error
	grpcConn, err := grpc.Dial(
		config.GetString("user_ms.address"),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatal("Can't connect to user microservice via grpc")
	}
	userManager = services.NewUserMSClient(grpcConn)
	defer grpcConn.Close()

	var nothing services.Nothing
	_, err = userManager.ClearDB(context.Background(), &nothing)
	if err != nil {
		t.Errorf(err.Error())
	}

	cases := []TestCaseAuth{
		TestCaseAuth{
			number:    1,
			email:     "test@test.com",
			nickname:  "boob",
			password1: "1209Qawsed",
			password2: "1209Qawsed",
		},
		TestCaseAuth{
			number:    2,
			email:     "1@test.com",
			nickname:  "asdasdsdasds",
			password1: "1209Qawsedbn",
			password2: "1209Qawsedbn",
		},
	}

	for _, user := range cases {

		form := url.Values{}
		form.Add("email", user.email)
		form.Add("nickname", user.nickname)
		form.Add("password", user.password1)
		form.Add("password2", user.password2)

		req := httptest.NewRequest(http.MethodPost, handlers.ApiRegister, nil)
		req.PostForm = form

		resp := httptest.NewRecorder()

		http.HandlerFunc(handlers.RegisterHandler).ServeHTTP(resp, req)

		if status := resp.Code; status == http.StatusInternalServerError {
			t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't make cookie or can't check validate",
				status)
		}
	}

	if err != nil {
		t.Error(err.Error())
	}
	for _, user := range cases {

		email := user.email
		password := user.password1

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
		}
	}
}

func TestAuthHandlerUnsuccessfulWrongFormsAndNotRegister(t *testing.T) {

	var userManager services.UserMSClient
	var err error
	grpcConn, err := grpc.Dial(
		config.GetString("user_ms.address"),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatal("Can't connect to user microservice via grpc")
	}
	userManager = services.NewUserMSClient(grpcConn)
	defer grpcConn.Close()

	var nothing services.Nothing
	_, err = userManager.ClearDB(context.Background(), &nothing)
	if err != nil {
		t.Errorf(err.Error())
	}

	cases := []TestCaseAuth{
		//TestCaseAuth{
		//	number:   1,
		//	email:    "test@1.com", //TODO with only numbers after @
		//	password1: "asdasdasdsadasdQ",
		//	error:    `{"email":"","password":"","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		//},
		TestCaseAuth{
			number:    2,
			email:     "test@@test.com",
			password1: "asdasdasdsadasdQ",
			error:     `{"email":"Неверно введён адрес электронной почты"}`,
		},
		TestCaseAuth{
			number:    3,
			email:     "te&st@test.com",
			password1: "asdasdasdsadasdQ",
			error:     `{"email":"Неверно введён адрес электронной почты"}`,
		},
		TestCaseAuth{
			number:    4,
			email:     "test@test.com",
			password1: "asdsdQ",
			error:     `{"password":"Пароль слишком короткий"}`,
		},
		TestCaseAuth{
			number:    5,
			email:     "teesttest.com",
			password1: "asdasdasdsadasdQ",
			error:     `{"email":"Неверно введён адрес электронной почты"}`,
		},
		//TestCaseAuth{
		//	number:   6,
		//	email:    "_____@test.com",
		//	password1: "asdasdasdsadasd",
		//	error:    `{"email":"","password":"","password2":"","nickname":"","avatar":"","error":["Пользователь с таким адресом электронной почты не зарегистрирован"]}`,
		//},
		TestCaseAuth{
			number:    7,
			email:     "teest@test.com",
			password1: "",
			error:     `{"password":"Пароль слишком короткий"}`,
		},
		TestCaseAuth{
			number:    8,
			email:     "",
			password1: "asdasdasdsadasd",
			error:     `{"email":"Неверно введён адрес электронной почты"}`,
		},
	}

	for _, item := range cases {
		email := item.email
		password := item.password1

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

	var userManager services.UserMSClient
	var err error
	grpcConn, err := grpc.Dial(
		config.GetString("user_ms.address"),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatal("Can't connect to user microservice via grpc")
	}
	userManager = services.NewUserMSClient(grpcConn)
	defer grpcConn.Close()

	var nothing services.Nothing
	_, err = userManager.ClearDB(context.Background(), &nothing)
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
		expected := `{"email":"Неверно введён адрес электронной почты","password":"Пароль слишком короткий"}`
		if response := resp.Body.String(); response != expected {
			t.Errorf("\nhandler returned wrong error response:\ngot %v\nwant %v;\n",
				response, expected)
		}
	}
}

func TestAuthDeleteHandlerSuccessful(t *testing.T) {

	var userManager services.UserMSClient
	var err error
	grpcConn, err := grpc.Dial(
		config.GetString("user_ms.address"),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatal("Can't connect to user microservice via grpc")
	}
	userManager = services.NewUserMSClient(grpcConn)
	defer grpcConn.Close()

	var nothing services.Nothing
	_, err = userManager.ClearDB(context.Background(), &nothing)
	if err != nil {
		t.Errorf(err.Error())
	}

	user := TestCaseAuth{
		number:    1,
		email:     "test@test.com",
		nickname:  "boob",
		password1: "1209Qawsed",
		password2: "1209Qawsed",
	}

	form := url.Values{}
	form.Add("email", user.email)
	form.Add("nickname", user.nickname)
	form.Add("password", user.password1)
	form.Add("password2", user.password2)

	req := httptest.NewRequest(http.MethodPost, handlers.ApiRegister, nil)
	req.PostForm = form

	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.RegisterHandler).ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't make cookie or can't check validate",
			status)
	}

	var user1 services.UserSignature
	user1.Password = user.password1
	user1.Email = user.email

	token, err := userManager.MakeToken(context.Background(), &user1)

	cookie := helpers.BuildSessionCookie(token)
	if err != nil {
		t.Errorf("MakeSession returned error: %s\n", err.Error())
		return
	}

	req = httptest.NewRequest(http.MethodDelete, handlers.ApiProfile, nil)
	req.AddCookie(&cookie)

	resp = httptest.NewRecorder()

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
