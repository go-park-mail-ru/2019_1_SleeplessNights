package TheGame

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/player"
	"github.com/gorilla/websocket"
	"math/rand"
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

func TestGameFacade_PlayByWebsocket(t *testing.T) {
	//Метод должен протестировать то, что по вызову метода game.PlayByWebsocket игрок попадает в game.in
	game.Close()//Там крутиться горутина, которая читает game.in, останавливаем её чтобы с ней не конкурировать
	game.in = make(chan player.Player)//Канал in был закрыт предыдущим методом, переопределяем его, чтобы протестировать
	uid := rand.Uint64()
	//TODO make valid websocket connection for this test to work
	game.PlayByWebsocket(&websocket.Conn{}, uid)
	newPlayer := <- game.in
	if newPlayer.UID() != uid {
		t.Error("game.PlayByWebsocket() violates UID")
	}
}

func TestGameFacade_StartBalance(t *testing.T) {
	
}
