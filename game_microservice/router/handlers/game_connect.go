package handlers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/game"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/messge"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("ChatMS")
	logger.SetLogLevel(logrus.TraceLevel)
}

func UpgradeWs(user *services.User, w http.ResponseWriter, r *http.Request) {
	logger.Info("Request entered UpgradeWs")
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
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

	gameInstance := game.GetInstance()
	go gameInstance.PlayByWebsocket(conn, uint64(user.Id))
}
