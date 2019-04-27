package handlers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/messge"

	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("ChatMS")
	logger.SetLogLevel(logrus.TraceLevel)
}

func EnterChat(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, http.Header{"Upgrade": []string{"websocket"}})
	logger.Info("Someone's connected to websocket chat")

	if err != nil {
		logger.Error(`micro service error in "EnterChat" during connection"`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = conn.WriteJSON(messge.Message{Title: "INFO", Payload: "you've been connected to Chat"})
	if err != nil {
		fmt.Println(err)
	}

	//TODO GET chatroom pointer, try to add user to chat as a new chat member
}
