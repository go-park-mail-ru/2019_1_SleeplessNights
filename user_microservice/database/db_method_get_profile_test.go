package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"testing"
)

func TestGetProfileSuccessful(t *testing.T) {

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

	oldProfile := services.Profile{
		User:    &oldUser,
		Matches: 0,
		WinRate: 0,
		Rating:  0,
	}

	_, err = database.GetInstance().AddUser(oldUser.Email, oldUser.Nickname, oldUser.AvatarPath, []byte{}, []byte{})
	if err != nil {
		t.Error(err.Error())
		return
	}

	profile, err := database.GetInstance().GetProfile(1)

	if profile.User.Nickname != oldUser.Nickname {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			profile.User.Nickname, oldUser.Nickname)
	}
	if profile.User.Id != oldUser.Id {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			profile.User.Id, oldUser.Id)
	}
	if profile.User.Email != oldUser.Email {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			profile.User.Email, oldUser.Email)
	}
	if profile.User.AvatarPath != oldUser.AvatarPath {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			profile.User.AvatarPath, oldUser.AvatarPath)
	}
	if profile.Matches != oldProfile.Matches {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			profile.User.AvatarPath, oldUser.AvatarPath)
	}
	if profile.WinRate != oldProfile.WinRate {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			profile.User.AvatarPath, oldUser.AvatarPath)
	}
	if profile.Rating != oldProfile.Rating {
		t.Errorf("DB returned wrong user:\ngot %v\nwant %v\n",
			profile.User.AvatarPath, oldUser.AvatarPath)
	}
}

func TestGetProfileUnsuccessful(t *testing.T) {

	expected := errors.DataBaseNoDataFound

	_, err := database.GetInstance().GetProfile(2)
	if err == nil {
		t.Errorf("DB didn't return any error")
		return
	} else if err.Error() != expected.Error() {
		t.Errorf("DB returned wrong error:\ngot %v\nwant %v\n",
			err.Error(), expected)
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}
