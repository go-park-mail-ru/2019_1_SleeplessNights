package user_manager_test

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/user_manager"
	"testing"
)

func TestClearDBSuccessful (t *testing.T) {
	_, err := user_manager.GetInstance().ClearDB(context.Background(), &services.Nothing{})
	if err != nil {
		t.Errorf("User_manager returned error: %v", err.Error())
		return
	}
}
