package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"testing"
)

func TestUpdateUserSuccessful(t *testing.T) {

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

	newUser := services.User{
		Id:         1,
		Email:      "new@test.com",
		Nickname:   "new",
		AvatarPath: "new.jpg",
	}

	err = database.GetInstance().UpdateUser(&newUser)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}

	user, err := database.GetInstance().GetUserByID(oldUser.Id)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
	if user.Nickname != newUser.Nickname || user.AvatarPath != newUser.AvatarPath {
		t.Errorf("DB returned wrong user:\ngot %v, %v\nwant %v, %v\n",
			user.Nickname, newUser.Nickname, user.AvatarPath, newUser.AvatarPath)
	}
}
