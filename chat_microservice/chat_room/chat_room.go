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

const (
	maxConnections        = 100
	limit          uint64 = 20
)

type chatRoom struct {
	maxConnections int64
	Id             uint64
	authorPool     map[uint64]*Author
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
		authorPool:     make(map[uint64]*Author),
	}
}

func GetInstance() *chatRoom {
	return chat
}

func (chat *chatRoom) Join(author Author) {
	logger.Info("User ", author.Nickname, "Joined room")

	chat.mx.Lock()
	chat.authorPool[author.Id] = &author
	wg := sync.WaitGroup{}
	wg.Add(1)
	logger.Info("Started Listening from User", author.Nickname)
	go func() {
		author.StartListen(chat.Id)
		wg.Done()
	}()
	chat.mx.Unlock()
	wg.Wait()
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
	AvatarPath string `json:"avatar_path"`
	Id         uint64 `json:"id"`
	Text       string `json:"text"`
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
				Nickname:   author.Nickname,
				AvatarPath: author.AvatarPath,
				Id:         author.Id,
				Text:       text,
			}

			bytes, err := json.Marshal(respMsg)
			if err != nil {
				logger.Error(err.Error())
			}

			err = database.GetInstance().PostMessage(respMsg.Id, roomId, bytes)
			if err != nil {
				logger.Error(err.Error())
			}

			for _, u := range chat.authorPool {
				err = u.Conn.WriteJSON(respMsg)
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
				err = author.Conn.WriteMessage(websocket.BinaryMessage, []byte(messages))
				if err != nil {
					logger.Error(err.Error())
				}
			}
		default:
			logger.Error("Message title not valid")
		}
	}
}
