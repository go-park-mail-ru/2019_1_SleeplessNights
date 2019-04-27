package handlers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/chat_room"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/messge"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/models"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

var logger *log.Logger

const (
	SendingRate = 0.1 // 1 /  SendingRate== Expected Value in seconds
)

func init() {
	logger = log.GetLogger("ChatMS")
	logger.SetLogLevel(logrus.TraceLevel)
}

func EnterChat(user models.User, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	logger.Info("Someone's connected to websocket chat, ID:%d", user.ID)

	if err != nil {
		logger.Error(`micro service error in "EnterChat" during connection"`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = conn.WriteJSON(messge.Message{Title: "INFO", Payload: "you've been connected to Chat"})
	if err != nil {
		fmt.Println(err)
	}
	isAuthorized := false
	if user.ID != 0 {
		isAuthorized = true
	}
	logger.Info("Users authStatus=%d ID=%d", isAuthorized, user.ID)
	//TODO GET chatroom pointer, try to add user to chat as a new chat member
	userId, err := database.GetInstance().UpdateUser(user.ID, user.Nickname, user.AvatarPath)
	if err != nil {
		logger.Error("Failed to get user in ChatConnect, from db.getI.UpdateUser ")
	}
	chat_room.GetInstance().Join(chat_room.Author{Conn: conn, Nickname: user.Nickname, AvatarPath: user.AvatarPath, Id: userId})
	//go MessageMux(user)
	go StartSendingTestMessages(conn)
}

func StartSendingTestMessages(conn *websocket.Conn) {
	for {
		sample := int64(rand.ExpFloat64() / SendingRate)
		time.Sleep(time.Duration(sample) * time.Second)
		err := conn.WriteJSON(messge.Message{Title: "INFO", Payload: `{"nickname":"IvanPetrov", "avatar_path:"/img/default_avatar.jpg", "text":"Hello, Grand Webmaster"}`})
		if err != nil {
			fmt.Println(err)
			break
		}
	}

}
