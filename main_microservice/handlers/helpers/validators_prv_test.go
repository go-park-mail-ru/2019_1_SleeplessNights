package helpers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	AvatarPrefix = "/main_microservice/static/img/"
)

type testCase struct {
	number   int
	email    string
	nickname string
	password string
}

func TestValidateEmailSuccessful(t *testing.T) {
	tests := []testCase{
		{
			number: 1,
			email:  "test@test.com",
		},
		{
			number: 2,
			email:  "tes1t@test.com",
		},
		{
			number: 3,
			email:  "hjghgj@test.com",
		},
		{
			number: 4,
			email:  "gg@gg.com",
		},
		{
			number: 5,
			email:  "h_h@test.com",
		},
		{
			number: 6,
			email:  "test@h-h.com",
		},
		{
			number: 7,
			email:  "qwerty@qwerty.ru",
		},
		{
			number: 8,
			email:  "a_k@k_.ruru",
		},
		{
			number: 9,
			email:  "a_s-g@jjhhj.hl",
		},
		{
			number: 10,
			email:  "g.2@com.com",
		},
	}

	for _, test := range tests {
		var requestErrors ErrorSet
		err := validateEmail(test.email, &requestErrors)
		if err != nil {
			t.Errorf("Number: %v\nValidator returned error: %v\n", test.number, err.Error())
		}
		if len(requestErrors) != 0 {
			t.Errorf("Number: %v\nValidator returned validation error: %v\n", test.number, requestErrors)
		}
	}
}

func TestValidateEmailUnsuccessful(t *testing.T) {
	tests := []testCase{
		{
			number: 1,
			email:  "@test.com",
		},
		{
			number: 2,
			email:  "tes1t@.com",
		},
		{
			number: 3,
			email:  "hjghgj@test.",
		},
		{
			number: 4,
			email:  "sdf!@gg.com",
		},
		{
			number: 5,
			email:  "ывапыва@test.com",
		},
		{
			number: 6,
			email:  "|sdf@h-h.com",
		},
		{
			number: 7,
			email:  "sdf:sf@qwerty.ru",
		},
		{
			number: 8,
			email:  "a_k@k_.r",
		},
		{
			number: 9,
			email:  "a",
		},
		{
			number: 10,
			email:  "g.2com.com",
		},
	}

	for _, test := range tests {
		var requestErrors ErrorSet
		err := validateEmail(test.email, &requestErrors)
		if err != nil {
			t.Errorf("Number: %v\nValidator returned error: %v\n", test.number, err.Error())
		}
		if len(requestErrors) == 0 {
			t.Errorf("Number: %v\nValidator didn't return validation error", test.number)
		}
	}
}

func TestValidatePasswordSuccessful(t *testing.T) {
	test := "sdfhjksdafhjkadsfhjk"

	var requestErrors ErrorSet
	err := validatePassword(test, &requestErrors)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	if len(requestErrors) != 0 {
		t.Errorf("Validator returned validation error: %v\n", requestErrors)
	}
}

func TestValidatePasswordUnsuccessful(t *testing.T) {
	test := "sss"

	var requestErrors ErrorSet
	err := validatePassword(test, &requestErrors)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	if len(requestErrors) == 0 {
		t.Errorf("Validator didn't return validation error")
	}
}

func TestValidateNicknameSuccessful(t *testing.T) {
	tests := []testCase{
		{
			number:   1,
			nickname: "boooooob",
		},
		{
			number:   2,
			nickname: "asd123_-",
		},
		{
			number:   3,
			nickname: "I-_-I",
		},
		{
			number:   4,
			nickname: "___KING___",
		},
		{
			number:   5,
			nickname: "123123123",
		},
	}

	for _, test := range tests {
		var requestErrors ErrorSet
		err := validateNickname(test.nickname, &requestErrors)
		if err != nil {
			t.Errorf("Number: %v\nValidator returned error: %v\n", test.number, err.Error())
		}
		if len(requestErrors) != 0 {
			t.Errorf("Number: %v\nValidator returned validation error: %v\n", test.number, requestErrors)
		}
	}
}

func TestValidateNicknameUnsuccessful(t *testing.T) {
	tests := []testCase{
		{
			number:   1,
			nickname: "прпр",
		},
		{
			number:   2,
			nickname: "()-()",
		},
		{
			number:   3,
			nickname: "asd",
		},
		{
			number:   4,
			nickname: "asdasdasdasdasdasd",
		},
		{
			number:   5,
			nickname: "?asd",
		},
		{
			number:   6,
			nickname: "!asd",
		},
		{
			number:   7,
			nickname: "\asd",
		},
	}

	for _, test := range tests {
		var requestErrors ErrorSet
		err := validateNickname(test.nickname, &requestErrors)
		if err != nil {
			t.Errorf("Number: %v\nValidator returned error: %v\n", test.number, err.Error())
		}
		if len(requestErrors) == 0 {
			t.Errorf("Number: %v\nValidator didn't return validation error", test.number)
		}
	}
}

func TestValidateAvatarSuccessful(t *testing.T) {

	img := "default_avatar.jpg"

	req, err := createReqWithAvatar(img)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.MultipartForm.File["avatar"][0].Header.Set("content-type", "image/jpeg")

	avatar := req.MultipartForm.File["avatar"][0]

	var requestErrors ErrorSet
	err = validateAvatar(avatar, &requestErrors)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	if len(requestErrors) != 0 {
		t.Errorf("Validator returned validation error: %v\n", requestErrors)
	}
}

func TestValidateAvatarUnsuccessfulEmptyFile(t *testing.T) {

	avatar := &multipart.FileHeader{
		Size: 0,
	}

	var requestErrors ErrorSet
	err := validateAvatar(avatar, &requestErrors)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	if len(requestErrors) == 0 {
		t.Errorf("Validator didn't return validation error")
	}
}

func TestValidateAvatarUnsuccessfulTooMuchSize(t *testing.T) {

	img := "test_avatar.jpeg"

	req, err := createReqWithAvatar(img)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.MultipartForm.File["avatar"][0].Header.Set("content-type", "image/jpeg")

	avatar := req.MultipartForm.File["avatar"][0]

	var requestErrors ErrorSet
	err = validateAvatar(avatar, &requestErrors)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	if len(requestErrors) == 0 {
		t.Errorf("Validator didn't return validation error")
	}
}

func TestValidateAvatarUnsuccessfulWrongType(t *testing.T) {

	img := "test_avatar.json"

	req, err := createReqWithAvatar(img)
	if err != nil {
		t.Errorf(err.Error())
	}
	req.MultipartForm.File["avatar"][0].Header.Set("content-type", "txt/json")

	avatar := req.MultipartForm.File["avatar"][0]

	var requestErrors ErrorSet
	err = validateAvatar(avatar, &requestErrors)
	if err != nil {
		t.Errorf("Validator returned error: %v\n", err.Error())
	}
	if len(requestErrors) == 0 {
		t.Errorf("Validator didn't return validation error")
	}
}

func createReqWithAvatar(img string) (req *http.Request, err error) {
	path := os.Getenv("BASEPATH") + AvatarPrefix + img

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("avatar", path)
	if err != nil {
		return
	}

	fh, err := os.Open(path)
	if err != nil {
		return
	}

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return
	}

	err = fh.Close()
	if err != nil {
		return
	}

	err = bodyWriter.Close()
	if err != nil {
		return
	}

	req = httptest.NewRequest(http.MethodPatch, "/", bodyBuf)
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	err = req.ParseMultipartForm(MaxPhotoSize)
	return
}
