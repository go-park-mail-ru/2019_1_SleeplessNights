package game

//TODO ADD TO PACKAGE:
//TODO - panic handling
//TODO - channels closing
//TODO - goroutines exits

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/player"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/room"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"sync"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Game")
	logger.SetLogLevel(logrus.Level(config.GetInt("game_ms.log_level")))
}

var (
	maxRooms                = config.GetInt("game_ms.pkg.game.max_rooms")
	maxPlayersInputQueueLen = config.GetInt("game_ms.pkg.game.player_input_queue_len")
)

//Game - это синглтон, механику их построения я подробнее объяснил, когда описывал PlayerFactory

//Игра - это фасад для всего микросервиса. Она должна предоставлять внешний интерфейс для всего, с чем мы можем
//начать игру (в базовом варианте - websocket соединение) и перерабатывать эти данные во внутренние абстракции -
//игроков, комнаты и т.д.

//Паттерн синглтон был выбран по той же причине, что и с фабрикой игроков, потому что зачем нам два не связанных
//между собой набора комнат игроков и т.д. в рамках одного приложения

var game *gameFacade
var in chan player.Player //Через этот канал игроки попадают из фабрики в комнаты
//in вынесен за пределы игры, чтобы улучшить масштабируемость на случай,
//если мы быдем поднимать несколько инстансов микросервиса игры

type gameFacade struct {
	maxRooms int //Макисмальное количество комнат в мапе, которое мы готовы поддерживать
	rooms    map[uint64]*room.Room
	idSource uint64
	mu       sync.Mutex
}

func init() {
	//Make Buffered channel, otherwise 'g.PlayByWebsocket' gets stuck,
	// awaiting for someone ('startBalance' in goroutine) to read from channel g.in
	in = make(chan player.Player, maxPlayersInputQueueLen)
	closer.Bind(func() {
		close(in) //Закрываем канал входа перед завершением работы программы
	})

	game = &gameFacade{
		maxRooms: maxRooms,
		//Сразу делаем мапу нужного размера, чтобы не тратить потом время на аллокации памяти
		//В принципе, если maxRooms большое и памяти жрёт много, то можно здесь поставить что-то типа
		//(maxRooms / 2)  или (maxRooms / 4)
		rooms:    make(map[uint64]*room.Room, maxRooms/4),
		idSource: 0,
	}
	go func() {
		game.startBalance() //Начинаем работу балансировщика
		err := recover()
		for err != nil {
			logger.Error("Unhandled panic came from startBalance:", err)
			game.startBalance()
			err = recover()
		}
	}()
}

func GetInstance() *gameFacade {
	return game
}
