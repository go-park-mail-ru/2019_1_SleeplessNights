package helpers_test

import (
	"bytes"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

type TestCaseReg struct {
	number    int
	email     string
	nickname  string
	password1 string
	password2 string
}

func TestValidateUpdateProfileRequestSuccessful(t *testing.T) {
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

	reqErr, err := helpers.ValidateUpdateProfileRequest(req)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	if len(reqErr) != 0 {
		t.Errorf("Validator returned validation error: %v\n", reqErr)
	}
}

func TestValidateUpdateProfileRequestUnsuccessful(t *testing.T) {

	img := "test_avatar.jpeg"
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

	newEmail := "g.2com.com"
	newNickname := "()-()"

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

	reqErr, err := helpers.ValidateUpdateProfileRequest(req)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	if len(reqErr) == 0 {
		t.Errorf("Validator didn't return validation error: %v\n", reqErr)
	}
}

func TestValidateRegisterRequestSuccessful(t *testing.T) {
	user := TestCaseReg{
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

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.PostForm = form
	err := req.ParseForm()
	if err != nil{
		t.Errorf("Parsing form returned error: %s\n", err.Error())
		return
	}

	reqErr, err := helpers.ValidateRegisterRequest(req)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	if len(reqErr) != 0 {
		t.Errorf("Validator returned validation error: %v\n", reqErr)
	}
}

func TestValidateRegisterRequestUnsuccessful(t *testing.T) {
	user := TestCaseReg{
		number:    1,
		email:     "g.2com.com",
		nickname:  "()-()",
		password1: "1209Qawsed",
		password2: "1209Qased",
	}

	form := url.Values{}
	form.Add("email", user.email)
	form.Add("nickname", user.nickname)
	form.Add("password", user.password1)
	form.Add("password2", user.password2)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.PostForm = form
	err := req.ParseForm()
	if err != nil{
		t.Errorf("Parsing form returned error: %s\n", err.Error())
		return
	}

	reqErr, err := helpers.ValidateRegisterRequest(req)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	if len(reqErr) == 0 {
		t.Errorf("Validator didn't return validation error: %v\n", reqErr)
	}
}