package room_manager_test

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/room_manager"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"testing"
)

func TestCreateRoomSuccessful(t *testing.T) {

	room := &services.RoomSettings{
		MaxConnections: uint64(2),
		Talkers:        []uint64{1, 2},
	}

	var ctx context.Context

	roomId, err := room_manager.GetInstance().CreateRoom(ctx, room)
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
		return
	}

	if roomId.Id != 2 {
		t.Errorf("CreateRoom returned wrong id:\ngot %v\nwnat %v",
			roomId.Id, 2)
	}
}

func TestCreateRoomUnsuccessful(t *testing.T) {

	room := &services.RoomSettings{
		MaxConnections: uint64(2),
	}

	var ctx context.Context

	_, err := room_manager.GetInstance().CreateRoom(ctx, room)
	if err == nil {
		t.Errorf("Didn't return error")
	}
}
