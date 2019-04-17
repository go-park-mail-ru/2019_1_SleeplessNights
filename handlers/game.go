package handlers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/gorilla/websocket"
	"net/http"
)

func UpgradeWs(user models.User, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, http.Header{"Upgrade": []string{"websocket"}})
	logger.Infof("player with ID = %d connected to socket", user.ID)
	if err != nil {
		logger.Error("Error during UpgradeWs", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = conn.WriteJSON(messge.Message{Title: "INFO", Payload: "you've been connected to server"})
	if err != nil {
		fmt.Println(err)
	}

	game := TheGame.GetInstance()
	game.PlayByWebsocket(conn, uint64(user.ID))

}
