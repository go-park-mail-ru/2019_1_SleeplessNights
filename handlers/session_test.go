package handlers_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/faker"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestAuthHandler(t *testing.T) {

	path := "/api/auth"

	faker.CreateFakeData(handlers.UserCounter)

	for  _, user := range models.Users{
		email := user.Email
		password := faker.FakeUserPassword

		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)

		req := httptest.NewRequest(http.MethodPost, path, nil)
		req.PostForm = form

		resp := httptest.NewRecorder()

		http.HandlerFunc(handlers.AuthHandler).ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code:\n got %v\n want %v;\n error: %s\n",
				status, http.StatusOK, resp.Body)
		}

		expected := `{}`
		if resp.Body.String() != expected {
			t.Errorf("handler returned unexpected body:\n got %v\n want %v\n",
				resp.Body.String(), expected)
		}
	}
}
