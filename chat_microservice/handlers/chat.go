package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/room_manager"
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
		roomId = room_manager.GlobalChatId
	} else {
		roomId, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			logger.Error(`Failed in getting query`, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(`Micro service error in "EnterChat" during connection`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Infof("Someone's connected to websocket chat, ID: %d", user.Id)

	if _, ok := room_manager.GetInstance().RoomsPool[roomId]; !ok {
		logger.Error(`Failed in finding room`)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	room := room_manager.GetInstance().RoomsPool[roomId]
	if uint64(len(room.TalkersPool)) == room.MaxConnections{
		logger.Error(`Failed because room is full`)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	var haveAccess bool
	for _, id := range room.AccessArray{
		if id == user.Id{
			haveAccess = true
		}
	}
	if !haveAccess{
		logger.Error(`Failed because user hasn't access to come into this room!`)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	isAuthorized := false
	if user.Id != 0 {
		isAuthorized = true
	}

	logger.Infof("Users authStatus - %v ID: %d", isAuthorized, user.Id)

	userId, err := database.GetInstance().UpdateTalker(user.Id, user.Nickname, user.AvatarPath)
	if err != nil {
		logger.Error("Failed to get user_manager in ChatConnect, from db.getI.UpdateTalker ")
	}
	room_manager.GetInstance().RoomsPool[roomId].Join(room_manager.Talker{
		Conn:       conn,
		Nickname:   user.Nickname,
		AvatarPath: user.AvatarPath,
		Id:         userId,
	})
}
