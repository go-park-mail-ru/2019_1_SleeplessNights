package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/room_manager"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

func EnterChat(user *services.User, w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("room")
	var roomId uint64
	var err error
	if str == "" {
		roomId = room_manager.GlobalChatId
	} else {
		roomId, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			logger.Error(`Failed in getting query:`, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if _, ok := room_manager.GetInstance().RoomsPool[roomId]; !ok {
		logger.Error(`Failed in finding room`)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	room := room_manager.GetInstance().RoomsPool[roomId]
	if uint64(len(room.TalkersPool)) == room.MaxConnections {
		logger.Error(`Failed because room is full`)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if room.Id != room_manager.GlobalChatId {
		var haveAccess bool
		for _, id := range room.AccessArray {
			if id == user.Id {
				haveAccess = true
			}
		}
		if !haveAccess {
			logger.Error(`Failed because user hasn't access to  enter this room!`)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
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

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(`Micro service error in "EnterChat" during connection`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Infof("Someone's connected to websocket chat, ID: %d", user.Id)

	room_manager.GetInstance().RoomsPool[roomId].Join(room_manager.Talker{
		Conn:       conn,
		Nickname:   user.Nickname,
		AvatarPath: user.AvatarPath,
		Id:         userId,
	})
}
