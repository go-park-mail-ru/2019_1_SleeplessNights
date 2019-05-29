package factory_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/player/factory"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/router/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"
	"github.com/gorilla/websocket"
	"net/http/httptest"
	"testing"
)

func TestBuildWebsocketPlayer (t *testing.T) {
	router := httptest.NewServer(middleware.MiddlewareAuth(handlers.UpgradeWs, false))

	d := websocket.Dialer{}

	_, _, err := d.Dial("ws://"+router.Listener.Addr().String(), nil)
	if err != nil {
		t.Errorf("Dial returned an error when trying to connect: %v", err.Error())
		return
	}
	router.Close()
}


func PlayerLogicFunc(id uint64, in, out *chan message.Message, args ...interface{}){}

func TestBuildChannelPlayer(t *testing.T) {
	pl := factory.GetInstance().BuildChannelPlayer(PlayerLogicFunc)
	pl.ID()
	pl.UID()
	pl.Subscribe()
}