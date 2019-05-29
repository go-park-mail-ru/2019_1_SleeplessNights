package handlers_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/router/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"
	"github.com/gorilla/websocket"
	"net/http/httptest"
	"testing"
)

func TestUpgradeWsSuccessful(t *testing.T) {

	router := httptest.NewServer(middleware.MiddlewareAuth(handlers.UpgradeWs, false))
	defer router.Close()

	d := websocket.Dialer{}

	_, _, err := d.Dial("ws://"+router.Listener.Addr().String(), nil)
	if err != nil {
		t.Errorf("Dial returned an error when trying to connect: %v", err.Error())
		return
	}
}
