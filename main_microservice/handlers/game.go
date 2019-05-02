package handlers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/messge"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/gorilla/websocket"
	"net/http"
)

func UpgradeWs(user services.User, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, http.Header{"Upgrade": []string{"websocket"}})
	logger.Infof("player with ID = %d connected to socket", user.Id)
	if err != nil {
		logger.Error("Error during UpgradeWs", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = conn.WriteJSON(messge.Message{Title: "INFO", Payload: "you've been connected to server"})
	if err != nil {
		fmt.Println(err)
	}

	game := game_microservice.GetInstance()
	go game.PlayByWebsocket(conn, uint64(user.Id))
}
