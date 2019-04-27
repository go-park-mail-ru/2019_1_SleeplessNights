package chat_room

import (
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
	mx         sync.Mutex
}

func init() {
	chat = &chatRoom{
		maxConnections: maxConnections,
	}
}

func GetInstance() *chatRoom {
	return chat
}

func (chat *chatRoom) Join(author Author) {
	logger.Info("User ", author.Nickname, "Joined room")
	chat.mx.Lock()
	chat.authorPool[author.Id] = author
	chat.mx.Unlock()
	logger.Info("Started Listening from User", author.Nickname)
	author.StartListen()
	chat.mx.Lock()
	logger.Info(" User", author.Nickname, "is Leaving Chat Room")
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

func (author *Author) StartListen() {
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
		case scrollTitle:
		default:

		}
	}
}
