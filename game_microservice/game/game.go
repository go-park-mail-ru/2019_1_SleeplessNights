package game

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/player/factory"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/room"
	"github.com/gorilla/websocket"
)

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
				for k, r := range g.rooms {
					if r.KillMePleaseFlag == true {
						r.CloseResponseRequestChannels()
						delete(g.rooms, k)
						logger.Info("Game, Room" + fmt.Sprint(k) + "was deleted from map")
					} else {
						if r.TryJoin(p) {
							logger.Info("Found Existing room_manager, player UID ", p.UID(), " added")
							roomFound = true
							break
						}

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
	//in <- factory.GetInstance().BuildChannelPlayer(jobToDo, args...)
	logger.Info("Player has been read from channel by balancer ")
}
