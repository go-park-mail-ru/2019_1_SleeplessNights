package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/game"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("GameMS-mux")
	logger.SetLogLevel(logrus.Level(config.GetInt("main_ms.log_level")))
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

	gameInstance := game.GetInstance()
	go gameInstance.PlayByWebsocket(conn, uint64(user.Id))
}
