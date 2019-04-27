package chat_room

import (
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/gorilla/websocket"
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
	AuthorPool map[uint64]Author
	//TODO add mutex
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
	//TODO join map
	//TODO start listen author
	//TODO kick author from map
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

