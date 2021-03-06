package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"testing"
)

func TestCleanerDBForTests(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	talkersArray := []uint64{1}

	_, err = database.GetInstance().AddRoom(talkersArray)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	rooms, err := database.GetInstance().GetRooms()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	if len(rooms) != 0 {
		t.Errorf("DB didn't cleaned up")
		return
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}
}