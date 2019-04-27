package chat_room

import (
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"golang.org/x/net/websocket"
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
}

func init() {
	chat = &chatRoom{
		maxConnections: maxConnections,
	}
}

func GetInstance() *chatRoom {
	return chat
}

type Author struct {
	Wc         *websocket.Conn
	Nickname   string
	AvatarPath string
	Id         uint64
}


