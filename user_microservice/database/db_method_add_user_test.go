package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/jackc/pgx"
	"testing"
)

func TestAddUserSuccessful(t *testing.T) {

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

	newUser, err := database.GetInstance().AddUser(oldUser.Email, oldUser.Nickname, oldUser.AvatarPath, []byte{}, []byte{})
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	if newUser.Email != oldUser.Email || newUser.Id != oldUser.Id {
		t.Errorf("DB returned wrong user:\ngot %v, %v\nwant %v, %v\n",
			newUser.Email, newUser.Id, oldUser.Email, oldUser.Id)
	}
}

func TestAddUserUnsuccessful(t *testing.T) {

	user := services.User{
		Email:      "test@test.com",
		Nickname:   "test",
		AvatarPath: "default_avatar.jpg",
	}

	expected := "23505"

	user, err := database.GetInstance().AddUser(user.Email, user.Nickname, user.AvatarPath, []byte{}, []byte{})
	if err == nil {
		t.Errorf("DB didn't return any error")
		return
	} else if err, ok := err.(pgx.PgError); ok && err.Code != expected {
		t.Errorf("DB returned wrong error:\ngot %v\nwant %v\n",
			err.Code, expected)
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
}
