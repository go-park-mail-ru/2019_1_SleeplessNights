package user_manager_test

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/user_manager"
	"testing"
)

func TestGetUserById(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	oldUser := services.NewUserData{
		Email:    "test@test.com",
		Nickname: "test",
	}


	_, err = user_manager.GetInstance().CreateUser(context.Background(), &oldUser)
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, err = user_manager.GetInstance().GetUserById(context.Background(), &services.UserId{Id:1})
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, err = user_manager.GetInstance().GetUserById(context.Background(), &services.UserId{Id:2})
	if err == nil {
		return
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}
