package user_manager_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/user_manager"
	"golang.org/x/net/context"
	"testing"
)

func TestCreateUserSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	oldUser := services.NewUserData{
		Email:    "test@test.com",
		Nickname: "test",
	}

	var ctx context.Context

	newUser, err := user_manager.GetInstance().CreateUser(ctx, &oldUser)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if newUser.Email != oldUser.Email || newUser.Nickname != oldUser.Nickname {
		t.Errorf("DB returned wrong user:\ngot %v, %v\nwant %v, %v\n",
			newUser.Email, newUser.Nickname, oldUser.Email, oldUser.Nickname)
	}
}

func TestCreateUserUnsuccessful(t *testing.T) {

	oldUser := services.NewUserData{
		Email:    "test@test.com",
		Nickname: "test",
	}

	var ctx context.Context
	expected := errors.DataBaseUniqueViolation
	_, err := user_manager.GetInstance().CreateUser(ctx, &oldUser)
	if err == nil {
		t.Errorf("DB didn't return any error")
		return
	} else if err != expected {
		t.Errorf("DB returned wrong error code:\ngot %v\nwant %v\n",
			err.Error(), expected.Error())
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}
