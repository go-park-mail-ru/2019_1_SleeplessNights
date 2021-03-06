package handlers_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestCaseReg struct {
	number    int
	email     string
	nickname  string
	password1 string
	password2 string
	error     string
}

func TestRegisterHandlerSuccessful(t *testing.T) {

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

	cases := []TestCaseReg{
		{
			number:    1,
			email:     "test@test.com",
			nickname:  "boob",
			password1: "1209Qawsed",
			password2: "1209Qawsed",
		},
		{
			number:    2,
			email:     "1@test.com",
			nickname:  "asdasdsdasds",
			password1: "1209Qawsedbn",
			password2: "1209Qawsedbn",
		},
		{
			number:    3,
			email:     "ghjklllll@asdasdasd.adasd",
			nickname:  "8667t68ff8",
			password1: "134fK34f34fed",
			password2: "134fK34f34fed",
		},
		{
			number:    4,
			email:     "acsdvs@tsdcsdcsdcsdcest.com",
			nickname:  "scklop",
			password1: "JNJasadasdasdasdNJ",
			password2: "JNJasadasdasdasdNJ",
		},
		{
			number:    5,
			email:     "iejuiwejfiuhihiufhiwufhiwh@test.ru",
			nickname:  "KKLKLLKKLKLKL",
			password1: "1209QaALLALALwsed",
			password2: "1209QaALLALALwsed",
		},
		{
			number:    6,
			email:     "aaaa@a.org",
			nickname:  "CDCDCDCDCDC",
			password1: "1092380912830912830128390183091839Qawsed",
			password2: "1092380912830912830128390183091839Qawsed",
		},
		{
			number:    7,
			email:     "zxczxczxc@zxcsdcdcewc.com",
			nickname:  "___jhghjgg",
			password1: "12938109238Hsdskdjhfksdhfkj",
			password2: "12938109238Hsdskdjhfksdhfkj",
		},
		{
			number:    8,
			email:     "test@asfsdfsdfsdfsdfdsfdsfsdf.com",
			nickname:  "090909090909",
			password1: "12ASDASDASDSADASDASD09Qawsed",
			password2: "12ASDASDASDSADASDASD09Qawsed",
		},
		{
			number:    9,
			email:     "test@yachoo.ru",
			nickname:  "booboijNo",
			password1: "booboijNo",
			password2: "booboijNo",
		},
		{
			number:    10,
			email:     "test2@test.com",
			nickname:  "120912129J",
			password1: "120912129J",
			password2: "120912129J",
		},
	}

	for _, item := range cases {

		form := url.Values{}
		form.Add("email", item.email)
		form.Add("nickname", item.nickname)
		form.Add("password", item.password1)
		form.Add("password2", item.password2)

		req := httptest.NewRequest(http.MethodPost, handlers.ApiRegister, nil)
		req.PostForm = form

		resp := httptest.NewRecorder()

		http.HandlerFunc(handlers.RegisterHandler).ServeHTTP(resp, req)

		if status := resp.Code; status == http.StatusInternalServerError {
			t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't make cookie or can't check validate",
				status)
		} else {
			if status := resp.Code; status != http.StatusOK {
				t.Errorf("\ntest %d.\nhandler returned wrong status code:\n got %v\n want %v\n",
					item.number, status, http.StatusOK)
			}
		}
	}
}

func TestRegisterHandlerUnsuccessfulWrongForms(t *testing.T) {

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

	cases := []TestCaseReg{
		{
			number:    1,
			email:     ".com", //TODO test@1.com
			nickname:  "boob",
			password1: "asdasdasdsadasd",
			password2: "asdasdasdsadasd",
			error:     `{"email":"Неверно введён адрес электронной почты"}`,
		},
		{
			number:    2,
			email:     "test.com",
			nickname:  "asdasdsdasds",
			password1: "1209Qawsedbn",
			password2: "1209Qawsedbn",
			error:     `{"email":"Неверно введён адрес электронной почты"}`,
		},
		{
			number:    3,
			email:     "----sdfsfds.adasd",
			nickname:  "8667t68ff8",
			password1: "134fK34f34fed",
			password2: "134fK34f34fed",
			error:     `{"email":"Неверно введён адрес электронной почты"}`,
		},
		{
			number:    4,
			email:     "iejihiufhiwufhiwh@test.ru",
			nickname:  "KKLKLLKKLKLKL",
			password1: "12d",
			password2: "12d",
			error:     `{"password":"Пароль слишком короткий"}`,
		},
		{
			number:    5,
			email:     "aasaa@a.org",
			nickname:  "CDCDCDCDCDC",
			password1: "109238091839Qawsed",
			password2: "1092380912830912830128390183091839Qawsed",
			error:     `{"password2":"Пароли не совпадают"}`,
		},
		{
			number:    6,
			email:     "zxczxczxc@zxdcdcewc.com",
			nickname:  "",
			password1: "12938109238Hsdskdjhfksdhfkj",
			password2: "12938109238Hsdskdjhfksdhfkj",
			error:     `{"nickname":"Никнейм не может быть короче 4 символов"}`,
		},
		{
			number:    7,
			email:     "test@asfdfsdfsdfdsfdsfsdf.com",
			nickname:  "090909090909",
			password1: "",
			password2: "",
			error:     `{"password":"Пароль слишком короткий"}`,
		},
		{
			number:    8,
			email:     "",
			nickname:  "",
			password1: "",
			password2: "1",
			error:     `{"email":"Неверно введён адрес электронной почты","password":"Пароль слишком короткий","password2":"Пароли не совпадают","nickname":"Никнейм не может быть короче 4 символов"}`,
		},
		{
			number:    9,
			email:     "",
			nickname:  "120912129Jsdfsdfsdfsdfsdfsdfsdfsdfsdf",
			password1: "12129J",
			password2: "120912129J",
			error:     `{"email":"Неверно введён адрес электронной почты","password2":"Пароли не совпадают","nickname":"Никнейм не может быть длиннее 16 символов"}`,
		},
	}

	for _, item := range cases {

		form := url.Values{}
		form.Add("email", item.email)
		form.Add("nickname", item.nickname)
		form.Add("password", item.password1)
		form.Add("password2", item.password2)

		req := httptest.NewRequest(http.MethodPost, handlers.ApiRegister, nil)
		req.PostForm = form

		resp := httptest.NewRecorder()

		http.HandlerFunc(handlers.RegisterHandler).ServeHTTP(resp, req)

		if status := resp.Code; status == http.StatusInternalServerError {
			t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't make cookie or can't check validate",
				status)
		} else {
			if status := resp.Code; status != http.StatusBadRequest {
				t.Errorf("\ntest %d.\nhandler returned wrong status code:\ngot %v\nwant %v\n",
					item.number, status, http.StatusBadRequest)
			}

			if resp.Body.String() != item.error {
				t.Errorf("\ntest %d.\nhandler returned wrong error response:\ngot %v\nwant %v\n",
					item.number, resp.Body.String(), item.error)
			}
		}
	}
}
