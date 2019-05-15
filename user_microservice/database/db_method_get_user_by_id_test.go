package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/jackc/pgx"
	"testing"
)

func TestGetUserByIdSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	oldUser := services.User{
		Id:         1,
		Email:      "test@test.com",
		Nickname:   "test",
		AvatarPath: "default_avatar.jpg",
	}

	_, err = database.GetInstance().AddUser(oldUser.Email, oldUser.Nickname, oldUser.AvatarPath, []byte{}, []byte{})
	if err != nil {
		t.Error(err.Error())
		return
	}

	newUser, err := database.GetInstance().GetUserByID(oldUser.Id)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	if newUser.Id != oldUser.Id {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			newUser.Id, oldUser.Id)
	}
	if newUser.Email != oldUser.Email {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			newUser.Email, oldUser.Email)
	}
	if newUser.Nickname != oldUser.Nickname {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			newUser.Nickname, oldUser.Nickname)
	}
	if newUser.AvatarPath != oldUser.AvatarPath {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			newUser.AvatarPath, oldUser.AvatarPath)
	}
}

func TestGetUserViaIdUnsuccessful(t *testing.T) {

	expected := "P0002"

	_, err := database.GetInstance().GetUserByID(100)
	if err == nil {
		t.Errorf("DB didn't return error")
	} else if err, ok := err.(pgx.PgError); ok && err.Code != expected {
		t.Errorf("DB returned wrong error:\ngot %v\nwant %v\n",
			err.Code, expected)
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}
