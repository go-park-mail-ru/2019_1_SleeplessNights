package room

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	maxConnections        = 100
	limit          uint64 = 20
)

var chat *room

var logger *log.Logger

func init() {
	logger = log.GetLogger("Room")
	logger.SetLogLevel(logrus.TraceLevel)
}

type room struct {
	maxConnections int64
	Id             uint64
	usersPool      map[uint64]*User
	mx             sync.Mutex
}

func init() {
	id, err := database.GetInstance().AddRoom(nil)
	if err != nil {
		logger.Error("Chat_room init", err)
	}
	chat = &room{
		Id:             id,
		maxConnections: maxConnections,
		usersPool:      make(map[uint64]*User),
	}
}

func GetInstance() *room {
	return chat
}

func (chat *room) Join(user User) {
	logger.Info("User ", user.Nickname, "Joined room")

	chat.mx.Lock()
	chat.usersPool[user.Id] = &user
	wg := sync.WaitGroup{}
	wg.Add(1)
	logger.Info("Started Listening from User", user.Nickname)
	go func() {
		user.StartListen(chat.Id)
		wg.Done()
	}()
	chat.mx.Unlock()
	wg.Wait()
	chat.mx.Lock()
	logger.Info(" User", user.Nickname, "is Leaving Chat Room")
	delete(chat.usersPool, user.Id)
	chat.mx.Unlock()
}

type User struct {
	Conn       *websocket.Conn
	Nickname   string
	AvatarPath string
	Id         uint64
}

type Message struct {
	Title   string      `json:"title"`
	Payload interface{} `json:"payload"`
}

const (
	postTitle   = "POST"
	scrollTitle = "SCROLL"
)

type PostPayload struct {
	Text string `json:"text,omitempty"`
}

type ScrollPayload struct {
	Since uint64 `json:"since"`
}

type ResponseMessage struct {
	Nickname   string `json:"nickname"`
	AvatarPath string `json:"avatarPath"`
	Id         uint64 `json:"id"`
	Text       string `json:"text"`
}

func (us *User) StartListen(roomId uint64) {

	var msg Message
	for {
		err := us.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				logger.Infof("Player %d closed the connection", us.Id)
				return
			}
		}
		logger.Info("Got Message from connection", msg)

		switch msg.Title {
		case postTitle:
			st, ok := msg.Payload.(map[string]interface{})

			logger.Info(st)
			if !ok {
				logger.Error("Something wrong with msg.Payload.(Post)")
			}
			text, ok := st["text"].(string)
			logger.Info(text)
			if !ok {
				logger.Error(`st[text] error`)
			}

			//TODO switch for payload types

			respMsg := ResponseMessage{
				Nickname:   us.Nickname,
				AvatarPath: us.AvatarPath,
				Id:         us.Id,
				Text:       text,
			}

			bytes, err := json.Marshal(respMsg)
			if err != nil {
				logger.Error(err.Error())
			}

			err = database.GetInstance().AddMessage(respMsg.Id, roomId, bytes)
			if err != nil {
				logger.Error(err.Error())
			}

			for _, user := range chat.usersPool {
				err = user.Conn.WriteJSON(respMsg)
				if err != nil {
					logger.Error(err.Error())
				}
			}
		case scrollTitle:
			{
				st, ok := msg.Payload.(map[string]interface{})
				logger.Info(st)
				if !ok {
					logger.Error("Something wrong with msg.Payload.(ScrollPayload)")
				}
				since, ok := st["since"].(float64)
				logger.Info(since)
				if !ok {
					logger.Error(`st[since] error`)
				}

				messages, err := database.GetInstance().GetMessages(roomId, uint64(since), limit)
				if err != nil {
					logger.Error(err.Error())
				}
				logger.Info(messages)
				err = us.Conn.WriteMessage(websocket.BinaryMessage, []byte(messages))
				if err != nil {
					logger.Error(err.Error())
				}
			}
		default:
			logger.Error("Message title not valid")
		}
	}
}
