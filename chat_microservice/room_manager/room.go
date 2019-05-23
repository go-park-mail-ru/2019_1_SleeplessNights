package room_manager

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/gorilla/websocket"
	"sync"
)

func (r *room) Join(user Talker) {
	logger.Info("User ", user.Nickname, "Joined room_manager")
	r.mx.Lock()
	r.UsersPool[user.Id] = &user
	wg := sync.WaitGroup{}
	wg.Add(1)
	logger.Info("Started Listening from User", user.Nickname)
	go func() {
		user.StartListen(r.id)
		wg.Done()
	}()
	r.mx.Unlock()
	wg.Wait()
	r.mx.Lock()
	logger.Info(" User", user.Nickname, "is Leaving Chat Room")
	delete(r.UsersPool, user.Id)
	r.mx.Unlock()
}

func (t *Talker) StartListen(roomId uint64) {
	var msg message
	for {
		err := t.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				logger.Infof("Talker %d closed the connection", t.Id)
				return
			}
		}
		logger.Info("Got Message from connection", msg)

		//TODO switch for payload types

		switch msg.Title {
		case postTitle:
			respMsg := responseMessage{
				Nickname:   t.Nickname,
				AvatarPath: t.AvatarPath,
				Id:         t.Id,
				Text:       msg.Payload.Text,
			}

			bytes, err := json.Marshal(respMsg)
			if err != nil {
				logger.Error(err.Error())
			}

			err = database.GetInstance().AddMessage(respMsg.Id, roomId, bytes)
			if err != nil {
				logger.Error(err.Error())
			}

			logger.Debugf("Len of user pool: %d", len(chat.RoomsPool[roomId].UsersPool))
			for _, user  := range chat.RoomsPool[roomId].UsersPool {
				err = user.Conn.WriteJSON(respMsg)
				if err != nil {
					logger.Error(err.Error())
				}
			}
		case scrollTitle:
			logger.Debug(msg.Payload.Since)
			logger.Debug(roomId)
			logger.Debug(limit)
			messages, err := database.GetInstance().GetMessages(roomId, uint64(msg.Payload.Since), limit)
			if err != nil {
				logger.Error(err.Error())
			}
			logger.Debug(messages)

			err = t.Conn.WriteMessage(websocket.BinaryMessage, []byte(messages))
			if err != nil {
				logger.Error(err.Error())
			}
		default:
			logger.Error("Message title not valid")
		}
	}
}
