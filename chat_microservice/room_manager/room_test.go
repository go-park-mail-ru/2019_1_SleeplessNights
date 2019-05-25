package room_manager_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJoinSuccessful(t *testing.T) {
	router := httptest.NewServer(middleware.MiddlewareAuth(handlers.EnterChat, false))
	defer router.Close()

	d := websocket.Dialer{}

	_, resp, err := d.Dial("ws://"+router.Listener.Addr().String(), nil)
	if err != nil {
		t.Errorf("Dial returned an error when trying to connect: %v", err.Error())
		return
	}

	if status := resp.StatusCode; status != http.StatusSwitchingProtocols {
		t.Errorf("handler returned wrong status code:\ngot %v\nwnat %v",
			status, http.StatusSwitchingProtocols)
	}
}

func TestStartListenSuccessful_POST(t *testing.T) {
	router := httptest.NewServer(middleware.MiddlewareAuth(handlers.EnterChat, false))
	defer router.Close()

	d := websocket.Dialer{}

	c, resp, err := d.Dial("ws://"+router.Listener.Addr().String(), nil)
	if err != nil {
		t.Errorf("Dial returned an error when trying to connect: %v", err.Error())
		return
	}

	if status := resp.StatusCode; status != http.StatusSwitchingProtocols {
		t.Errorf("handler returned wrong status code:\ngot %v\nwnat %v",
			status, http.StatusSwitchingProtocols)
	}

	err = c.WriteMessage(websocket.BinaryMessage, []byte(`{"title":"POST","payload":{"since":"Hi"}}`))
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartListenSuccessful_SCROLL(t *testing.T) {
	router := httptest.NewServer(middleware.MiddlewareAuth(handlers.EnterChat, false))
	defer router.Close()

	d := websocket.Dialer{}

	c, resp, err := d.Dial("ws://"+router.Listener.Addr().String(), nil)
	if err != nil {
		t.Errorf("Dial returned an error when trying to connect: %v", err.Error())
		return
	}

	if status := resp.StatusCode; status != http.StatusSwitchingProtocols {
		t.Errorf("handler returned wrong status code:\ngot %v\nwnat %v",
			status, http.StatusSwitchingProtocols)
	}

	err = c.WriteMessage(websocket.BinaryMessage, []byte(`{"title":"SCROLL","payload":{"since":50}}`))
	if err != nil {
		t.Fatal(err)
	}
}

func TestStartListenUnsuccessful_WrongTitle(t *testing.T) {
	router := httptest.NewServer(middleware.MiddlewareAuth(handlers.EnterChat, false))
	defer router.Close()

	d := websocket.Dialer{}

	c, resp, err := d.Dial("ws://"+router.Listener.Addr().String(), nil)
	if err != nil {
		t.Errorf("Dial returned an error when trying to connect: %v", err.Error())
		return
	}

	if status := resp.StatusCode; status != http.StatusSwitchingProtocols {
		t.Errorf("handler returned wrong status code:\ngot %v\nwnat %v",
			status, http.StatusSwitchingProtocols)
	}

	err = c.WriteMessage(websocket.BinaryMessage, []byte(`{"title":"t","payload":{"since":50}}`))
	if err != nil {
		t.Fatal(err)
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
}
