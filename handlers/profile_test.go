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

func TestProfileHandlerSuccessfulWithCreateFakeData(t *testing.T) {

	faker.CreateFakeData(handlers.UserCounter)

	for _, user := range models.Users{
		cookie, err := helpers.MakeSession(user)
		if err != nil {
			t.Errorf("\nMakeSession returned error: %s\n", err)
			return
		}

		req := httptest.NewRequest(http.MethodGet, handlers.ApiProfile, nil)
		req.AddCookie(&cookie)

		resp := httptest.NewRecorder()

		http.HandlerFunc(handlers.ProfileHandler).ServeHTTP(resp, req)
		if status := resp.Code; status == http.StatusInternalServerError {
			t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't Marshal 'user' into json\n",
				status)
		} else {
			if status := resp.Code; status != http.StatusOK {
				t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
					status, http.StatusOK)
			}

			expected :=
				fmt.Sprintf("{\"email\":\"%s\",\"won\":%d,\"lost\":%d,\"play_time\":%d,\"nickname\":\"%s\",\"avatar_path\":\"%s\"}",
					user.Email, user.Won, user.Lost, user.PlayTime, user.Nickname, user.AvatarPath)
			if resp.Body.String() != expected {
				t.Errorf("\nhandler returned unexpected body:\ngot %v\nwant %v\n",
					resp.Body.String(), expected)
			}
		}
	}
}

func TestProfileHandlerUnsuccessfulWithoutCookie(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, handlers.ApiProfile, nil)

	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.ProfileHandler).ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't Marshal 'user' into json\n",
			status)
	} else {
		if status := resp.Code; status != http.StatusUnauthorized {
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
				status, http.StatusUnauthorized)
		}

		expected := `{}`
		if resp.Body.String() != expected {
			t.Errorf("\nhandler returned unexpected body:\ngot %v\nwant %v\n",
				resp.Body.String(), expected)
		}
	}
}

func TestProfileHandlerUnsuccessfulWithWrongCookie(t *testing.T) {

	user := models.User{
		ID: 1000,
	}

	cookie, err := helpers.MakeSession(user)
	if err != nil {
		t.Errorf("MakeSession returned error: %s", err)
		return
	}

	req := httptest.NewRequest(http.MethodGet, handlers.ApiProfile, nil)
	req.AddCookie(&cookie)

	resp := httptest.NewRecorder()

	http.HandlerFunc(handlers.ProfileHandler).ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't Marshal 'user' into json\n",
			status)
	} else {
		if status := resp.Code; status != http.StatusUnauthorized {
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
				status, http.StatusUnauthorized)
		}

		expected := `{}`
		if resp.Body.String() != expected {
			t.Errorf("\nhandler returned unexpected body:\ngot %v\nwant %v\n",
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
