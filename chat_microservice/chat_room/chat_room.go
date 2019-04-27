package chat_room

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/gorilla/websocket"
	"sync"
)

var chat *chatRoom

var logger *log.Logger

func init() {
	logger = log.GetLogger("DB")
}

const (
	maxConnections = 100
)

type chatRoom struct {
	maxConnections int64
	authorPool map[uint64]Author
	mx sync.Mutex
}

func init() {
	chat = &chatRoom{
		maxConnections: maxConnections,
	}
}

func GetInstance() *chatRoom {
	return chat
}

func (chat *chatRoom)Join(author Author) {
	chat.mx.Lock()
	chat.authorPool[author.Id] = author
	chat.mx.Unlock()
	chat.authorPool[author.Id].StartListen()
	chat.mx.Lock()
	delete(chat.authorPool, author.Id)
	chat.mx.Unlock()
}

type Author struct {
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
	postTitle   = "post"
	scrollTitle = "scroll"
)

type PostPayload struct {
	Text string `json:"text,omitempty"`
}

type ScrollPayload struct {
	Since uint64 `json:"since"`
}

type ResponseMessage struct {
	Nickname string `json:"nickname"`
	AvatarPath string `json:"avatar_path"`
	Id uint64 `json:"id"`
	Text string `json:"text"`
}

func (author *Author)StartListen(roomId uint64) {
	var msg Message
	for {
		err := author.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				logger.Infof("Player %d closed the connection", author.Id)
				return
			}
		}
		logger.Info("Got from connection", msg)

		switch msg.Title {
		case postTitle:
			postPayload, ok := msg.Payload.(PostPayload)
			if !ok {
				logger.Error("Invalid payload came with post title, got error while spelling")
				break
			}

			//TODO switch for payload types

			respMsg := ResponseMessage{
				//TODO make struct
			}

			database.GetInstance().PostMessage(author.Id, roomId,)//TODO add message to DB

			//TODO for each user in room send message

		case scrollTitle:
		default:

		}
	}
}

