package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"sync"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Room")
	logger.SetLogLevel(logrus.Level(config.GetInt("chat_ms.log_level")))
}

const (
	postTitle   = "POST"
	scrollTitle = "SCROLL"
)

var (
	maxConnections = int64(config.GetInt("chat_ms.pkg.room.max_connections"))
	limit          = uint64(config.GetInt("chat_ms.pkg.room.msg_limit"))
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
