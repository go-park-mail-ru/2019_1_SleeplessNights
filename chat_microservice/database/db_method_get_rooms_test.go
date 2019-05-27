package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"testing"
)

func TestGetRoomsSuccessful(t *testing.T) {

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

	rooms, err := database.GetInstance().GetRooms()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	if len(rooms) != 1 {
		t.Errorf("DB didn't return rooms")
		return
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}
}
