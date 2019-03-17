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

	faker.CreateFakeData(1)
	email := models.Users[models.UserKeyPairs[1]].Email
	password := faker.FakeUserPassword

	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)

	path := "/api/auth"

	req := httptest.NewRequest(http.MethodPost, path, nil)
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.AuthHandler).ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v; error: %s",
			status, http.StatusOK, resp.Body)
	}

	expected := `{}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Body.String(), expected)
	}
}
