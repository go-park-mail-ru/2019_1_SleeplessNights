package TheGame

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/player"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/player/factory"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/room"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/gorilla/websocket"
	"sync"
)

//Game - это синглтон, механику их построения я подробнее объяснил, когда описывал PlayerFactory

//Игра - это фасад для всего микросервиса. Она должна предоставлять внешний интерфейс для всего, с чем мы можем
//начать игру (в базовом варианте - websocket соединение) и перерабатывать эти данные во внутренние абстракции -
//игроков, комнаты и т.д.

//Паттерн синглтон был выбран по той же причине, что и с фабрикой игроков, потому что зачем нам два не связанных
//между собой набора комнат игроков и т.д. в рамках одного приложения

const (
	maxRooms = 100
)

var game *gameFacade

type gameFacade struct {
	in       chan player.Player //Через этот канал игроки попадают из фабрики в комнаты
	maxRooms int                //Макисмальное количество комнат в мапе, которое мы готовы поддерживать
	rooms    map[uint64]room.Room
	idSource uint64
	mu       sync.Mutex
}

func init() {
	game = &gameFacade{
		maxRooms: maxRooms,
		//Сразу делаем мапу нужного размера, чтобы не тратить потом время на аллокации памяти
		//В принципе, если maxRooms большое и памяти жрёт много, то можно здесь поставить что-то типа
		//(maxRooms / 2)  или (maxRooms / 4)
		rooms:    make(map[uint64]room.Room, maxRooms/4),
		idSource: 0,
	}
	go game.startBalance() //Начинаем работу балансировщика
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
	//TODO develop
LOOP:
	for {
		select {
		case p := <-g.in:
			logger.Info.Println("player %s joined the game", p.ID())
			//Search for  room aplayer can join
			roomsCounter := 0
			for _, v := range g.rooms {
				if v.IsAvailable {
					if v.TryJoin(p) {
						break LOOP
					}
				}
				roomsCounter++
			}

			if roomsCounter != maxRooms {
				g.mu.Lock()
				g.idSource += 1
				roomId := g.idSource
				g.rooms[roomId] = room.Room{}
				if g.rooms[roomId].TryJoin(p) {
					break LOOP
				}
				g.mu.Unlock()
				logger.Info.Println("Successfully created Room with id %d", roomId)
				logger.Info.Println("Player with id %d added to room %d", p.ID(), roomId)
			} else {
				logger.Info.Println("Couldn't create new room for player %d. Number of rooms has reached the limit ", p.ID())

				//Send Message to Player
			}
		}
	}
}

func (g *gameFacade) PlayByWebsocket(conn *websocket.Conn) {
	//Начинаем игру по вебсокет соединению
	g.in <- factory.GetInstance().BuildWebsocketPlayer(conn) //Соьственно, всё изи
}
