package room

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
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

type roomManager struct {
	RoomsPool map[uint64]*room
}

var chat *roomManager

func init() {
	roomIds, err := database.GetInstance().GetRoomsIds()
	if err != nil {
		logger.Error("Chat_room init", err)
	}

	roomsPool := make(map[uint64]*room)
	for _, r := range roomIds {
		var room = &room{
			id:             r,
			maxConnections: maxConnections,
			usersPool:      make(map[uint64]*Talker),
		}
		roomsPool[room.id] = room
	}

	chat = &roomManager{
		RoomsPool: roomsPool,
	}
}

func GetInstance() *roomManager {
	return chat
}

func CreateRoom(id uint64) (r *room){
	r = &room{
		id:             id,
		maxConnections: maxConnections,
		usersPool:      make(map[uint64]*Talker),
	}
	return
}
