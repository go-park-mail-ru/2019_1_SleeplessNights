package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"testing"
)

func TestGetMessagesSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	uid := uint64(1)
	nickname := "test"
	avatarPath := "default_avatar.jpg"

	tId, err := database.GetInstance().UpdateTalker(uid, nickname, avatarPath)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	talkersArray := []uint64{1}

	rId, err := database.GetInstance().AddRoom(talkersArray)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	payload := `{"nickname":"Guest","avatarPath":"default_avatar.jpg","id":1,"text":"WOWOWOW"}`

	err = database.GetInstance().AddMessage(tId, rId, []byte(payload))
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	messages, err := database.GetInstance().GetMessages(rId, 100, 1)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	if len(messages) == 0 {
		t.Errorf("DB didn't return messages")
		return
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}
}

