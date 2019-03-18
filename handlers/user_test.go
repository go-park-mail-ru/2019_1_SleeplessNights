package handlers_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type TestCase struct {
	number    int
	email     string
	nickname  string
	password1 string
	password2 string
}

func TestRegisterHandler(t *testing.T) {

	path := "/api/register"

	cases := []TestCase{
		TestCase{
			number: 1,
			email:     "test@test.com",
			nickname:  "boob",
			password1: "1209Qawsed",
			password2: "1209Qawsed",
		},
		TestCase{
			number: 2,
			email:     "1@test.com",
			nickname:  "asdasdsdasds",
			password1: "1209Qawsedbn",
			password2: "1209Qawsedbn",
		},
		TestCase{
			number: 3,
			email:     "ghjklllll@asdasdasd.adasd",
			nickname:  "8667t68ff8",
			password1: "134fK34f34fed",
			password2: "134fK34f34fed",
		},
		TestCase{
			number: 4,
			email:     "acsdvs@tsdcsdcsdcsdcest.com",
			nickname:  "scklop",
			password1: "JNJasadasdasdasdNJ",
			password2: "JNJasadasdasdasdNJ",
		},
		TestCase{
			number: 5,
			email:     "iejuiwejfiuhihiufhiwufhiwh@test.ru",
			nickname:  "KKLKLLKKLKLKL",
			password1: "1209QaALLALALwsed",
			password2: "1209QaALLALALwsed",
		},
		TestCase{
			number: 6,
			email:     "aaaa@a.org",
			nickname:  "CDCDCDCDCDC",
			password1: "1092380912830912830128390183091839Qawsed",
			password2: "1092380912830912830128390183091839Qawsed",
		},
		TestCase{
			number: 7,
			email:     "zxczxczxc@zxcsdcdcewc.com",
			nickname:  "___jhghjgg",
			password1: "12938109238Hsdskdjhfksdhfkj",
			password2: "12938109238Hsdskdjhfksdhfkj",
		},
		TestCase{
			number: 8,
			email:     "test@asfsdfsdfsdfsdfdsfdsfsdf.com",
			nickname:  "090909090909",
			password1: "12ASDASDASDSADASDASD09Qawsed",
			password2: "12ASDASDASDSADASDASD09Qawsed",
		},
		TestCase{
			number: 9,
			email:     "test@yachoo.ru",
			nickname:  "booboijNo",
			password1: "booboijNo",
			password2: "booboijNo",
		},
		TestCase{
			number: 10,
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

		req := httptest.NewRequest(http.MethodPost, path, nil)
		req.PostForm = form

		resp := httptest.NewRecorder()

		http.HandlerFunc(handlers.RegisterHandler).ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("Test %d.\nhandler returned wrong status code:\n got %v\n want %v\n",
				item.number, status, http.StatusOK)
		}

		expected := `{}`
		if resp.Body.String() != expected {
			t.Errorf("Test %d.\nhandler returned unexpected body:\n got %v\n want %v\n",
				item.number, resp.Body.String(), expected)
		}
	}
}
