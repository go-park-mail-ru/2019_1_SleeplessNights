package game

//TODO ADD TO PACKAGE:
//TODO - panic handling
//TODO - channels closing
//TODO - goroutines exits

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/player"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/player/factory"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/room"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"sync"
)

//Game - это синглтон, механику их построения я подробнее объяснил, когда описывал PlayerFactory

//Игра - это фасад для всего микросервиса. Она должна предоставлять внешний интерфейс для всего, с чем мы можем
//начать игру (в базовом варианте - websocket соединение) и перерабатывать эти данные во внутренние абстракции -
//игроков, комнаты и т.д.

//Паттерн синглтон был выбран по той же причине, что и с фабрикой игроков, потому что зачем нам два не связанных
//между собой набора комнат игроков и т.д. в рамках одного приложения

var logger *log.Logger

func init() {
	logger = log.GetLogger("Game")
	logger.SetLogLevel(logrus.Level(config.GetInt("game_ms.log_level")))
}

var (
	maxRooms                = config.GetInt("game_ms.pkg.game.max_rooms")
	maxPlayersInputQueueLen = config.GetInt("game_ms.pkg.game.player_input_queue_len")
)

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

func (g *gameFacade) startBalance() {
	//Процедура запускает балансировщик игроков
	//Она в цикле читает канал in и, когда туда приходит игрок, распределяет его в комнату
	//Если свободных комнат нет, то создаёт новую

	//Стратегий распределения может быть много, в самом простом варианте мы идём в цикле по мапе комнат
	//и ищем в каждой комнате свободное место, если дошли до конца и не нашли, то создаём свою комнату
	//и занимаем место в ней, а если достигнут maxRooms, то заново входим в цикл
	logger.Trace("StartBalance started")

	for p := range in {

		fmt.Println("Got value from channel")
		logger.Info("Got new Player from channel g.in")
		go func() {
			err := p.Send(message.Message{Title: message.RoomSearching})
			if err != nil {
				logger.Warning("Failed to notify player with UID", p.UID())
			}

			logger.Info("goroutine Started:", "player UID", p.UID(), " looking for space room_manager")
			roomFound := false
			for !roomFound {
				roomsCounter := 0
				roomFound = false
				//Search for room_manager a player can join
				for _, v := range g.rooms {
					if v.TryJoin(p) {
						logger.Info("Found Existing room_manager, player UID ", p.UID(), " added")
						roomFound = true
						break
					}
					roomsCounter++
				}
				g.mu.Lock()
				var roomId uint64
				if roomsCounter != maxRooms && !roomFound {
					g.idSource += 1
					roomId = g.idSource
					g.rooms[roomId] = &room.Room{}
					roomFound = (g.rooms[roomId]).TryJoin(p)
				}
				g.mu.Unlock()
				if roomFound {
					logger.Info("Successfully found Room with id", roomId)
					logger.Info("Player with UID", p.UID(), "added to room_manager", roomId)
				} else {
					logger.Error("Failed to join just created Room with id", roomId)
				}
			}
		}()
	}
}

func (g *gameFacade) PlayByWebsocket(conn *websocket.Conn, uid uint64) {
	logger.Info("PlayByWebsocket Got new Connection")
	//Начинаем игру по вебсокет соединению
	in <- factory.GetInstance().BuildWebsocketPlayer(conn, uid) //Собственно, всё изи
	logger.Info("Player has been read from channel by balancer ")
}

func (g *gameFacade) PlayByChannels(jobToDo factory.ChannelPlayerLogic, args ...interface{}) {
	g.in <- factory.GetInstance().BuildChannelPlayer(jobToDo, args...)
	logger.Info("Player has been read from channel by balancer ")
}
