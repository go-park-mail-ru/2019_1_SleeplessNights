package handlers_test

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/faker"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProfileHandler(t *testing.T) {

	path := "/api/profile"

	faker.CreateFakeData(handlers.UserCounter)

	for _, user := range models.Users{
		cookie, err := helpers.MakeSession(user)
		if err != nil {
			t.Errorf("MakeSession returned error: %s", err)
			return
		}

		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.AddCookie(&cookie)

		resp := httptest.NewRecorder()

		http.HandlerFunc(handlers.ProfileHandler).ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code:\n got %v\n want %v\n",
				status, http.StatusOK)
		}

		expected :=
			fmt.Sprintf("{\"email\":\"%s\",\"won\":%d,\"lost\":%d,\"play_time\":%d,\"nickname\":\"%s\",\"avatar_path\":\"%s\"}",
				user.Email, user.Won, user.Lost, user.PlayTime, user.Nickname, user.AvatarPath)
		if resp.Body.String() != expected {
			t.Errorf("handler returned unexpected body:\n got %v\n want %v\n",
				resp.Body.String(), expected)
		}
	}
}

func TestProfileUpdateHandler(t *testing.T) {
	//req, err := http.NewRequest("PATCH", "/profile", nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//resp := httptest.NewRecorder()
	//handler := http.HandlerFunc(handlers.ProfileUpdateHandler)
	//
	//handler.ServeHTTP(resp, req)
	//
	//if status := resp.Code; status != http.StatusOK {
	//	t.Errorf(WrongStatus+": got %v want %v",
	//		status, http.StatusOK)
	//}
	//
	//expected := `{}` //TODO expected
	//if resp.Body.String() != expected {
	//	t.Errorf(UnexpectedBody+": got %v want %v",
	//		resp.Body.String(), expected)
	//}
}
