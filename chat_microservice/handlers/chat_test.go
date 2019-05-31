package handlers_test

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/room_manager"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnterChatSuccessful_PrivateChat(t *testing.T) {

	room := &services.RoomSettings{
		MaxConnections: 2,
		Talkers:        []uint64{0, 1},
	}

	var ctx context.Context

	_, err := room_manager.GetInstance().CreateRoom(ctx, room)
	if err != nil {
		t.Errorf("Room_manager returned an error when trying to create: %v", err.Error())
		return
	}

	router := httptest.NewServer(middleware.MiddlewareAuth(handlers.EnterChat, false))
	defer router.Close()

	d := websocket.Dialer{}

	_, resp, err := d.Dial("ws://"+router.Listener.Addr().String()+"?room=2", nil)
	if err != nil {
		t.Errorf("Dial returned an error when trying to connect: %v", err.Error())
		return
	}

	if status := resp.StatusCode; status != http.StatusSwitchingProtocols {
		t.Errorf("handler returned wrong status code:\ngot %v\nwnat %v",
			status, http.StatusSwitchingProtocols)
	}
}

func TestEnterChatSuccessful_GlobalChat(t *testing.T) {

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

func TestEnterChatUnsuccessful_WrongPath(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/api/chat", nil)
	req.URL.Query().Add("room", "h")
	resp := httptest.NewRecorder()

	http.Handler(middleware.MiddlewareAuth(handlers.EnterChat, false)).ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code:\ngot %v\nwnat %v",
			status, http.StatusBadRequest)
	}
}

func TestEnterChatUnsuccessful_WrongRoom(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/api/chat?room=3", nil)
	resp := httptest.NewRecorder()

	http.Handler(middleware.MiddlewareAuth(handlers.EnterChat, false)).ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code:\ngot %v\nwnat %v",
			status, http.StatusNotFound)
	}
}

func TestEnterChatUnsuccessful_FullRoom(t *testing.T) {

	router := httptest.NewServer(middleware.MiddlewareAuth(handlers.EnterChat, false))
	defer router.Close()

	d := websocket.Dialer{}

	_, _, err := d.Dial("ws://"+router.Listener.Addr().String()+"?room=2", nil)
	if err != nil {
		t.Errorf("Dial returned an error when trying to connect: %v", err.Error())
		return
	}

	req := httptest.NewRequest(http.MethodGet, "/api/chat?room=2", nil)
	resp := httptest.NewRecorder()

	http.Handler(middleware.MiddlewareAuth(handlers.EnterChat, false)).ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code:\ngot %v\nwnat %v",
			status, http.StatusServiceUnavailable)
	}
}

func TestEnterChatUnsuccessful_HaveNoAccess(t *testing.T) {

	room := &services.RoomSettings{
		MaxConnections: 2,
		Talkers:        []uint64{1, 2},
	}

	var ctx context.Context

	_, err := room_manager.GetInstance().CreateRoom(ctx, room)
	if err != nil {
		t.Errorf("Room_manager returned an error when trying to create: %v", err.Error())
		return
	}

	req := httptest.NewRequest(http.MethodGet, "/api/chat?room=3", nil)
	resp := httptest.NewRecorder()

	http.Handler(middleware.MiddlewareAuth(handlers.EnterChat, false)).ServeHTTP(resp, req)

	if status := resp.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code:\ngot %v\nwnat %v",
			status, http.StatusServiceUnavailable)
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf("DB returned error: %v", err.Error())
	}
}
