package handlers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/room"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Handlers")
	logger.SetLogLevel(logrus.TraceLevel)
}

const (
	SendingRate = 0.1 // 1 /  SendingRate== Expected Value in seconds
)

func EnterChat(user *services.User, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	logger.Infof("Someone's connected to websocket chat, ID: %d", user.Id)
	if err != nil {
		logger.Error(`Micro service error in "EnterChat" during connection"`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isAuthorized := false
	if user.Id != 0 {
		isAuthorized = true
	}

	logger.Infof("Users authStatus - %v ID: %d", isAuthorized, user.Id)

	//TODO GET chatroom pointer, try to add user_manager to chat as a new chat member

	userId, err := database.GetInstance().UpdateUser(user.Id, user.Nickname, user.AvatarPath)
	if err != nil {
		logger.Error("Failed to get user_manager in ChatConnect, from db.getI.UpdateUser ")
	}
	room.GetInstance().Join(room.User{
		Conn:       conn,
		Nickname:   user.Nickname,
		AvatarPath: user.AvatarPath,
		Id:         userId})
	//go MessageMux(user_manager)
	//go StartSendingTestMessages(conn)
}

func StartSendingTestMessages(conn *websocket.Conn) {
	for {
		sample := int64(rand.ExpFloat64() / SendingRate)
		time.Sleep(time.Duration(sample) * time.Second)
		err := conn.WriteJSON(message.Message{Title: "INFO", Payload: `{"nickname":"IvanPetrov", "avatar_path":"/default_avatar.jpg", "text":"Hello, Grand Webmaster"}`})
		if err != nil {
			fmt.Println(err)
			break
		}
	}

}
