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
	Wc         *websocket.Conn
	Nickname   string
	AvatarPath string
	Id         uint64
}

func (author *Author)StartListen() {
	//TODO handle messages
}

