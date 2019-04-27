package chat_room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/gorilla/websocket"
	"sync"
)

var chat *chatRoom

var logger *log.Logger

const (
	maxConnections        = 100
	limit          uint64 = 5
)

type chatRoom struct {
	maxConnections int64
	Id             uint64
	authorPool     map[uint64]Author
	mx             sync.Mutex
}

func init() {
	logger = log.GetLogger("DB")
	id, err := database.GetInstance().CreateRoom(nil)
	if err != nil {
		logger.Error("Chat_room init", err)
	}
	chat = &chatRoom{
		Id:             id,
		maxConnections: maxConnections,
		authorPool:     make(map[uint64]Author),
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
	author.StartListen(chat.Id)
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

func (author *Author) StartListen(roomId uint64) {

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
			{
				/*logger.Info(msg.Payload)
				sp, ok := msg.Payload.(ScrollPayload)
				logger.Info(sp)
				if !ok {
					logger.Error("Something wrong with msg.Payload.(ScrollPayload)")
				}*/

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
				err = author.Conn.WriteJSON([]byte(messages))
				if err != nil {
					logger.Error(err.Error())
				}
			}
		default:
			logger.Error("Message title not valid")
		}
	}
}
