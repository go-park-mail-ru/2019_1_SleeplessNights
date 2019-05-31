package handlers_test

import (
	"bytes"
	"context"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/user_manager"
	"google.golang.org/grpc"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestProfileHandlerSuccessfulWithCreateFakeData(t *testing.T) {

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

	user := services.NewUserData{
		Email:    "test@test.com",
		Nickname: "test",
		Password: "1209wawsed",
	}

	_, err = user_manager.GetInstance().CreateUser(context.Background(), &user)
	if err != nil {
		t.Error(err.Error())
		return
	}

	sig := &services.UserSignature{
		Email: "test@test.com",
		Password: "1209wawsed",
	}

	token, err := userManager.MakeToken(context.Background(), sig)
	if err != nil {
		t.Errorf(err.Error())
	}

	cookie := helpers.BuildSessionCookie(token)
	if err != nil {
		t.Errorf("\nMakeSession returned error: %s\n", err)
		return
	}

	req := httptest.NewRequest(http.MethodGet, handlers.ApiProfile, nil)
	req.AddCookie(&cookie)

	resp := httptest.NewRecorder()

	middleware.MiddlewareAuth(handlers.ProfileHandler, true).ServeHTTP(resp, req)
	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't Marshal 'user_manager' into json\n",
			status)
	} else {
		if status := resp.Code; status != http.StatusOK {
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
				status, http.StatusOK)
		}
	}

	_, err = userManager.ClearDB(context.Background(), &nothing)
	if err != nil {
		t.Errorf(err.Error())
	}
}


func TestProfileHandlerUnsuccessfulWithoutCookie(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, handlers.ApiProfile, nil)

	resp := httptest.NewRecorder()

	middleware.MiddlewareAuth(handlers.ProfileHandler, true).ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("\nhandler returned wrong status code: %v\nhandler can't write into responce or can't Marshal 'user_manager' into json\n",
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

	user := services.NewUserData{
		Email:    "test@test.com",
		Nickname: "test",
		Password: "1209wawsed",
	}

	_, err = user_manager.GetInstance().CreateUser(context.Background(), &user)
	if err != nil {
		t.Error(err.Error())
		return
	}

	sig := &services.UserSignature{
		Email: "test@test.com",
		Password: "1209wawsed",
	}

	token, err := userManager.MakeToken(context.Background(), sig)
	if err != nil {
		t.Errorf(err.Error())
	}

	cookie := helpers.BuildSessionCookie(token)
	if err != nil {
		t.Errorf("\nMakeSession returned error: %s\n", err)
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

	_, err = userManager.ClearDB(context.Background(), &nothing)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestProfileUpdateHandlerUnsuccessful_WrongRequest(t *testing.T) {

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

	user := services.NewUserData{
		Email:    "test@test.com",
		Nickname: "test",
		Password: "1209wawsed",
	}

	_, err = user_manager.GetInstance().CreateUser(context.Background(), &user)
	if err != nil {
		t.Error(err.Error())
		return
	}

	sig := &services.UserSignature{
		Email: "test@test.com",
		Password: "1209wawsed",
	}

	token, err := userManager.MakeToken(context.Background(), sig)
	if err != nil {
		t.Errorf(err.Error())
	}

	cookie := helpers.BuildSessionCookie(token)
	if err != nil {
		t.Errorf("\nMakeSession returned error: %s\n", err)
		return
	}

	req := httptest.NewRequest(http.MethodPatch, handlers.ApiProfile, nil)
	if err != nil {
		if err != nil {
			t.Errorf("Parsed returned error: %s\n", err.Error())
			return
		}
	}
	req.AddCookie(&cookie)

	resp := httptest.NewRecorder()

	middleware.MiddlewareAuth(handlers.ProfileUpdateHandler, true).ServeHTTP(resp, req)
	if status := resp.Code; status != http.StatusBadRequest {
		t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\nbody %v\n",
			status, http.StatusBadRequest, resp.Body)
	}

	_, err = userManager.ClearDB(context.Background(), &nothing)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestProfileUpdateHandlerUnsuccessful_NoValid(t *testing.T) {

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

	user := services.NewUserData{
		Email:    "test@test.com",
		Nickname: "test",
		Password: "1209wawsed",
	}

	_, err = user_manager.GetInstance().CreateUser(context.Background(), &user)
	if err != nil {
		t.Error(err.Error())
		return
	}

	sig := &services.UserSignature{
		Email: "test@test.com",
		Password: "1209wawsed",
	}

	token, err := userManager.MakeToken(context.Background(), sig)
	if err != nil {
		t.Errorf(err.Error())
	}

	cookie := helpers.BuildSessionCookie(token)
	if err != nil {
		t.Errorf("\nMakeSession returned error: %s\n", err)
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

	newEmail := "secondmail.com"
	newNickname := "secondaaaaaaaaaaaaaaaaaa"

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

	if status := resp.Code; status != http.StatusBadRequest {
		t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\nbody %v\n",
			status, http.StatusBadRequest, resp.Body)
	}

	_, err = userManager.ClearDB(context.Background(), &nothing)
	if err != nil {
		t.Errorf(err.Error())
	}
}
