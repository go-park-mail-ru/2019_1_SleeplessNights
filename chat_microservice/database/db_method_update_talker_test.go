package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"testing"
)

func TestUpdateTalkerSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	uid := uint64(1)
	nickname := "test"
	avatarPath := "default_avatar.jpg"

	_, err = database.GetInstance().UpdateTalker(uid, nickname, avatarPath)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	_, err = database.GetInstance().UpdateTalker(uid, nickname, avatarPath)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}
}