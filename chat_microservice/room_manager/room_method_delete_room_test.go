package room_manager_test

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/room_manager"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"testing"
)

func TestDeleteRoomSuccessful(t *testing.T) {

	room := &services.RoomId{
		Id: uint64(2),
	}

	var ctx context.Context

	_, err := room_manager.GetInstance().DeleteRoom(ctx, room)
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
		t.Errorf("DB didn't delete error: %v", err.Error())
	}
}

func TestDeleteRoomUnsuccessful_DeleteGlobalChat(t *testing.T) {

	room := &services.RoomId{
		Id: uint64(1),
	}

	var ctx context.Context

	_, err := room_manager.GetInstance().DeleteRoom(ctx, room)
	if err == nil {
		t.Errorf("DB didn't return error")
		return
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
}