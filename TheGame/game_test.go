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

