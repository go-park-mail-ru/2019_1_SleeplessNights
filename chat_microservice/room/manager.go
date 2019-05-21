package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"sync"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Room")
	logger.SetLogLevel(logrus.TraceLevel)
}

const (
	postTitle   = "POST"
	scrollTitle = "SCROLL"
)

const (
	maxConnections        = 100
	limit          uint64 = 20
)

type room struct {
	maxConnections int64
	Id             uint64
	usersPool      map[uint64]*Talker
	mx             sync.Mutex
}

var chat *room

func init() {
	id, err := database.GetInstance().AddRoom(nil)
	if err != nil {
		logger.Error("Chat_room init", err)
	}
	chat = &room{
		Id:             id,
		maxConnections: maxConnections,
		usersPool:      make(map[uint64]*Talker),
	}
}

func GetInstance() *room {
	return chat
}
