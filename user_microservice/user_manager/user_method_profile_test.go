package user_manager_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/user_manager"
	"golang.org/x/net/context"
	"testing"
)

func TestGetProfileSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	oldUser := services.NewUserData{
		Email:    "test@test.com",
		Nickname: "test",
	}

	var ctx context.Context

	_, err = user_manager.GetInstance().CreateUser(ctx, &oldUser)
	if err != nil {
		t.Error(err.Error())
		return
	}

	user := services.User{
		Id:       1,
		Email:    "test@test.com",
		Nickname: "test",
	}
	profile, err := user_manager.GetInstance().GetProfile(ctx, &user)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if profile.User.Email != oldUser.Email || profile.User.Nickname != oldUser.Nickname {
		t.Errorf("DB returned wrong user:\ngot %v, %v\nwant %v, %v\n",
			profile.User.Email, profile.User.Nickname, oldUser.Email, oldUser.Nickname)
	}
}

func TestGetProfileUnsuccessful(t *testing.T) {

	var ctx context.Context

	user := services.User{
		Id:       1000,
		Email:    "test@test.com",
		Nickname: "test",
	}

	expected := errors.DataBaseNoDataFound

	_, err := user_manager.GetInstance().GetProfile(ctx, &user)
	if err == nil {
		t.Errorf("DB didn't return any error")
		return
	} else if err.Error() != expected.Error() {
		t.Errorf("DB returned wrong error:\ngot %v\nwant %v\n",
			err.Error(), expected)
	}

}

func TestUpdateProfileSuccessful(t *testing.T) {

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

	var ctx context.Context

	user, err := user_manager.GetInstance().UpdateProfile(ctx, &newUser)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if user.AvatarPath != newUser.AvatarPath || user.Nickname != newUser.Nickname {
		t.Errorf("DB returned wrong user:\ngot %v, %v\nwant %v, %v\n",
			user.AvatarPath, user.Nickname, newUser.AvatarPath, newUser.Nickname)
	}
}

