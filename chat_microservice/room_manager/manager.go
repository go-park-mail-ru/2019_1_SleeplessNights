package room_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/jackc/pgx"
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

const (
	GlobalChatId = uint64(1)
)

var (
	maxConnections = uint64(config.GetInt("chat_ms.pkg.room_manager.max_connections"))
	limit          = uint64(config.GetInt("chat_ms.pkg.room_manager.msg_limit"))
)

type roomManager struct {
	RoomsPool map[uint64]*room
	Mx        sync.Mutex
}

var chat *roomManager

func init() {
	roomIds, err := database.GetInstance().GetRoomsIds()
	if err != nil {
		logger.Error("Chat_room init", err.Error())
	}

	roomsPool := make(map[uint64]*room)
	for _, r := range roomIds {
		roomsPool[r] = createRoom(r, maxConnections)
	}

	if len(roomsPool) == 0 {
		roomsPool[GlobalChatId] = createRoom(GlobalChatId, maxConnections)
	}

	chat = &roomManager{
		RoomsPool: roomsPool,
	}
}

func GetInstance() *roomManager {
	return chat
}

func createRoom(id uint64, maxConn uint64) (r *room) {
	r = &room{
		id:             id,
		MaxConnections: maxConn,
		UsersPool:      make(map[uint64]*Talker),
	}
	return
}

const (
	nodataFound         = "P0002"
	foreignKeyViolation = "23503"
)

func handlerError(pgError pgx.PgError) (err error) {
	switch pgError.Code {
	case foreignKeyViolation:
		err = errors.DataBaseForeignKeyViolation
	case nodataFound:
		err = errors.DataBaseNoDataFound
	default:
		err = pgError
	}
	return
}
