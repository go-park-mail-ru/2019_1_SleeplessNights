package tests

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfileHandler(t *testing.T) {

	var tempPassword []byte
	var tempSalt []byte
	testUser := models.User{
		ID:         10,
		Email:      "test@test.com",
		Password:   tempPassword,
		Salt:       tempSalt,
		Won:        1,
		Lost:       2,
		PlayTime:   10,
		Nickname:   "bob",
		AvatarPath: "/static/img/default_avatar.jpg",
	}

	models.Users[testUser.Email] = testUser

	cookie, err := router.MakeSession(testUser)
	if err != nil {
		t.Errorf("MakeSession returned error: %s", err)
		return
	}

	req, err := http.NewRequest("GET", "/api/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&cookie)

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ProfileHandler)

	handler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"nickname":"bob","email":"test@test.com","won":1,"lost":2,"play_time":10,"avatar_path":"/static/img/default_avatar.jpg"}`
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Body.String(), expected)
	}
}

func TestProfileUpdateHandler(t *testing.T) {
	req, err := http.NewRequest("PATCH", "/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ProfileUpdateHandler)

	handler.ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{}` //TODO expected
	if resp.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp.Body.String(), expected)
	}
}
