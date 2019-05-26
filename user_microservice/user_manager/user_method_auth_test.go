package user_manager_test

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/user_manager"
	"testing"
)

func TestMakeTokenAndCheckTokenSuccessful (t *testing.T){

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}

	u := services.NewUserData{
		Email:    "test@test.com",
		Nickname: "test",
		Password: "1209qawsed",
	}

	var ctx context.Context

	_, err = user_manager.GetInstance().CreateUser(ctx, &u)
	if err != nil {
		t.Errorf("User_manager returned error: %v", err.Error())
		return
	}

	sig := &services.UserSignature{
		Email: "test@test.com",
		Password: "1209qawsed",
	}

	token, err := user_manager.GetInstance().MakeToken(ctx, sig)
	if err != nil {
		t.Errorf("User_manager returned error: %v", err.Error())
		return
	}

	_, err = user_manager.GetInstance().CheckToken(ctx, token)
	if err != nil {
		t.Errorf("User_manager returned error: %v", err.Error())
		return
	}
}

func TestMakeTokenAndCheckTokenUnsuccessful (t *testing.T){

	var ctx context.Context

	sig := &services.UserSignature{
		Email: "test@test.com",
		Password: "1209wawsed",
	}

	_, err := user_manager.GetInstance().MakeToken(ctx, sig)
	if err == nil {
		t.Errorf("User_manager didn't return error: %v", err.Error())
		return
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
}