package handlers_test

import (
	"bytes"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/faker"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/models"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestProfileHandlerSuccessfulWithCreateFakeData(t *testing.T) {

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
		cookie, err := helpers.BuildSessionCookie(user.ID)
		if err != nil {
			t.Errorf("\nMakeSession returned error: %s\n", err)
			return
		}

		req := httptest.NewRequest(http.MethodGet, handlers.ApiProfile, nil)
		req.AddCookie(&cookie)

		resp := httptest.NewRecorder()

		middleware.MiddlewareAuth(handlers.ProfileHandler, true).ServeHTTP(resp, req)
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

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	req := httptest.NewRequest(http.MethodGet, handlers.ApiProfile, nil)

	resp := httptest.NewRecorder()

	middleware.MiddlewareAuth(handlers.ProfileHandler, true).ServeHTTP(resp, req)

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

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	user := models.User{
		ID: 1000,
	}

	cookie, err := helpers.BuildSessionCookie(user.ID)
	if err != nil {
		t.Errorf("MakeSession returned error: %s", err)
		return
	}

	req := httptest.NewRequest(http.MethodGet, handlers.ApiProfile, nil)
	req.AddCookie(&cookie)

	resp := httptest.NewRecorder()

	middleware.MiddlewareAuth(handlers.ProfileHandler, true).ServeHTTP(resp, req)

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

func TestProfileUpdateHandlerSuccessful(t *testing.T) {

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

	img := "default_avatar.jpg"
	path := os.Getenv("BASEPATH") + handlers.AvatarPrefix + img

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("avatar", path)
	if err != nil {
		t.Errorf("CreatFormFile returned error: %s\n", err.Error())
		return
	}

	fh, err := os.Open(path)
	if err != nil {
		t.Errorf("Open file returned error: %s\n", err.Error())
		return
	}

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		t.Errorf("Copy file returned error: %s\n", err.Error())
		return
	}

	err = fh.Close()
	if err != nil {
		t.Errorf("Close file returned error: %s\n", err.Error())
		return
	}

	newEmail := "second@mail.com"
	newNickname := "second"

	err = bodyWriter.WriteField("email", newEmail)
	if err != nil {
		t.Errorf("Added field returned error: %s\n", err.Error())
		return
	}

	err = bodyWriter.WriteField("nickname", newNickname)
	if err != nil {
		t.Errorf("Added field returned error: %s\n", err.Error())
		return
	}

	err = bodyWriter.Close()
	if err != nil {
		t.Errorf("Close file returned error: %s\n", err.Error())
		return
	}

	req := httptest.NewRequest(http.MethodPatch, handlers.ApiProfile, bodyBuf)
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	err = req.ParseMultipartForm(helpers.MaxPhotoSize)
	if err != nil {
		if err != nil {
			t.Errorf("Parsed returned error: %s\n", err.Error())
			return
		}
	}

	req.MultipartForm.File["avatar"][0].Header.Set("content-type", "image/jpeg")
	req.AddCookie(&cookie)

	resp := httptest.NewRecorder()

	middleware.MiddlewareAuth(handlers.ProfileUpdateHandler, true).ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code: %v\n",
			status)
	} else {
		if status := resp.Code; status != http.StatusOK {
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\nbody %v\n",
				status, http.StatusOK, resp.Body)
		}
	}

	user, err = database.GetInstance().GetUserViaID(1)
	if err != nil {
		t.Error(err.Error())
	}

	if user.Nickname != newNickname {
		t.Errorf("\nDB returned wrong nickname:\ngot %v\nwant %v\n",
			user.Email, newNickname)
	}
}

func TestProfileUpdateHandlerUnsuccessfulWithoutCookie(t *testing.T) {

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

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	err = bodyWriter.Close()
	if err != nil {
		t.Errorf("Close file returned error: %s\n", err.Error())
		return
	}

	req := httptest.NewRequest(http.MethodPatch, handlers.ApiProfile, bodyBuf)
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	resp := httptest.NewRecorder()

	middleware.MiddlewareAuth(handlers.ProfileUpdateHandler, true).ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code: %v\n",
			status)
	} else {
		if status := resp.Code; status != http.StatusUnauthorized {
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
				status, http.StatusUnauthorized)
		}
	}
}

//func TestProfileUpdateHandlerUnsuccessfulWithWrongCookie(t *testing.T) {
//
//	user := models.User{
//		ID:         1,
//		Email:      "first@mail.com",
//		Nickname:   "first",
//		AvatarPath: "none",
//	}
//
//	cookie, err := auth_microservice.MakeSession(user)
//	if err != nil {
//		t.Errorf("MakeSession returned error: %s\n", err.Error())
//		return
//	}
//
//	bodyBuf := &bytes.Buffer{}
//	bodyWriter := multipart.NewWriter(bodyBuf)
//
//	err = bodyWriter.Close()
//	if err != nil {
//		t.Errorf("Close file returned error: %s\n", err.Error())
//		return
//	}
//
//	req := httptest.NewRequest(http.MethodPatch, handlers.ApiProfile, bodyBuf)
//	req.AddCookie(&cookie)
//
//	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
//
//	resp := httptest.NewRecorder()
//
//	router.MiddlewareAuth(handlers.ProfileUpdateHandler).ServeHTTP(resp, req)
//
//	if status := resp.Code; status == http.StatusInternalServerError {
//		t.Errorf("\nhandler returned wrong status code: %v\n",
//			status)
//	} else {
//		if status := resp.Code; status != http.StatusUnauthorized {
//			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
//				status, http.StatusUnauthorized)
//		}
//	}
//}

func TestProfileUpdateHandlerUnsuccessfulWithoutMultipartForm(t *testing.T) {

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
		t.Errorf("AddUser returned error: %s\n", err.Error())
		return
	}

	cookie, err := helpers.BuildSessionCookie(user.ID)
	if err != nil {
		t.Errorf("MakeSession returned error: %s\n", err.Error())
		return
	}

	req := httptest.NewRequest(http.MethodPatch, handlers.ApiProfile, nil)
	req.AddCookie(&cookie)

	resp := httptest.NewRecorder()

	middleware.MiddlewareAuth(handlers.ProfileUpdateHandler, true).ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code: %v\n",
			status)
	} else {
		if status := resp.Code; status != http.StatusBadRequest {
			//Авторизация теперь вынесена в мидлвары и проводится раньше валидации
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
				status, http.StatusBadRequest)
		}
	}
}
