package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"testing"
)

func TestCleanerDBForTests(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	user := services.User{
		Email:      "test@test.com",
		Nickname:   "test",
		AvatarPath: "default_avatar.jpg",
	}

	user, err = database.GetInstance().AddUser(user.Email, user.Nickname, user.AvatarPath, []byte{}, []byte{})
	if err != nil {
		t.Error(err.Error())
		return
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	length, err := database.GetInstance().GetLenUsers()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	if length != 0 {
		t.Errorf("DB didn't cleaned up")
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}