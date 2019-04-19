package TheGame

import (
	"reflect"
	"testing"
)

func TestGetInstance(t *testing.T) {
	if reflect.TypeOf(GetInstance()) != reflect.PtrTo(reflect.TypeOf(gameFacade{})) {
		t.Errorf("GetInstance method returns value with wrong type: got %s, whant %s",
			reflect.TypeOf(reflect.TypeOf(GetInstance())),
			reflect.PtrTo(reflect.TypeOf(gameFacade{})))
	}
}

//func TestGameFacade_PlayByWebsocket(t *testing.T) {
//	game = &gameFacade{
//		maxRooms: 1,
//		rooms:    make(map[uint64]*room.Room, maxRooms),
//		idSource: 0,
//		in:       make(chan player.Player, 1),
//	}
//
//	uid := rand.Uint64()
//	game.PlayByWebsocket(&websocket.Conn{}, uid)
//	newPlayer := <-game.in
//
//	if newPlayer.UID() != uid {
//		t.Errorf("PlayByWebsocket method violates uid: got %d, whant %d", newPlayer.UID(), uid)
//	}
//}
