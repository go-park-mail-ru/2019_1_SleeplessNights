package room

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/gorilla/websocket"
	"sync"
)

func (r *room) Join(user Talker) {
	logger.Info("User ", user.Nickname, "Joined room")
	r.mx.Lock()
	r.usersPool[user.Id] = &user
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
	delete(r.usersPool, user.Id)
	if len(r.usersPool) == 0 && r.id != 1 {
		chat.Mx.Lock()
		delete(chat.RoomsPool, r.id)
		chat.Mx.Unlock()
	}
	r.mx.Unlock()
}

func (us *Talker) StartListen(roomId uint64) {
	var msg Message
	for {
		err := us.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				logger.Infof("Talker %d closed the connection", us.Id)
				return
			}
		}
		logger.Info("Got Message from connection", msg)

		//TODO switch for payload types

		switch msg.Title {
		case postTitle:
			respMsg := ResponseMessage{
				Nickname:   us.Nickname,
				AvatarPath: us.AvatarPath,
				Id:         us.Id,
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

			logger.Debugf("Len of user pool: %d", len(chat.RoomsPool[roomId].usersPool))
			for _, user := range chat.RoomsPool[roomId].usersPool {
				err = user.Conn.WriteJSON(respMsg)
				if err != nil {
					logger.Error(err.Error())
				}
			}
		case scrollTitle:
			logger.Debug(msg.Payload.Since)
			messages, err := database.GetInstance().GetMessages(roomId, uint64(msg.Payload.Since), limit)
			if err != nil {
				logger.Error(err.Error())
			}
			logger.Debug(messages)

			err = us.Conn.WriteMessage(websocket.BinaryMessage, []byte(messages))
			if err != nil {
				logger.Error(err.Error())
			}
		default:
			logger.Error("Message title not valid")
		}
	}
}
