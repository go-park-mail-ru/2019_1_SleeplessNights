package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/room"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Handlers")
	logger.SetLogLevel(logrus.Level(config.GetInt("chat_ms.log_level")))
}

func EnterChat(user *services.User, w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	str := r.URL.Query().Get("room")
	var roomId uint64
	var err error
	if str == "" {
		roomId = 1
	} else {
		roomId, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			logger.Error(`Failed in getting query"`, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	logger.Infof("Someone's connected to websocket chat, ID: %d", user.Id)
	if err != nil {
		logger.Error(`Micro service error in "EnterChat" during connection"`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, ok := room.GetInstance().RoomsPool[roomId]; !ok {
		r := room.CreateRoom(roomId)
		room.GetInstance().RoomsPool[roomId] = r
	}

	isAuthorized := false
	if user.Id != 0 {
		isAuthorized = true
	}

	logger.Infof("Users authStatus - %v ID: %d", isAuthorized, user.Id)

	//TODO GET chatroom pointer, try to add user_manager to chat as a new chat member

	userId, err := database.GetInstance().UpdateTalker(user.Id, user.Nickname, user.AvatarPath)
	if err != nil {
		logger.Error("Failed to get user_manager in ChatConnect, from db.getI.UpdateTalker ")
	}
	err = room.GetInstance().RoomsPool[roomId].Join(room.Talker{
		Conn:       conn,
		Nickname:   user.Nickname,
		AvatarPath: user.AvatarPath,
		Id:         userId,
	})

	if err != nil {
		logger.Error(`Failed in deleting room from db`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
