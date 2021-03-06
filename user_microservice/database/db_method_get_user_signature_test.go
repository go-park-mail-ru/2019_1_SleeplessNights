package database_test

import (
	"bytes"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/jackc/pgx"
	"testing"
)

func TestGetUserSignatureSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}

	oldUser := services.User{
		Id:         1,
		Email:      "test@test.com",
		Nickname:   "test",
		AvatarPath: "default_avatar.jpg",
	}

	_, err = database.GetInstance().AddUser(oldUser.Email, oldUser.Nickname, oldUser.AvatarPath, []byte{}, []byte{})
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	var zero []byte
	id, password, salt, err := database.GetInstance().GetUserSignature(oldUser.Email)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	if id != oldUser.Id {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			id, oldUser.Id)
	}
	if !bytes.Equal(password, zero){
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			password, zero)
	}
	if !bytes.Equal(salt, zero) {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			salt, zero)
	}
}

func TestGetUserSignatureUnsuccessful(t *testing.T) {

	var userEmail = "test88@test.com"

	expected := "P0002"

	_, _, _, err := database.GetInstance().GetUserSignature(userEmail)
	if err == nil {
		t.Errorf("DB didn't return error")
	} else if err, ok := err.(pgx.PgError); ok && err.Code != expected {
		t.Errorf("DB returned wrong error:\ngot %v\nwant %v\n",
			err.Code, expected)
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
}