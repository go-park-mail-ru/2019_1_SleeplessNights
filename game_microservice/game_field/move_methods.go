package game_field

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/event"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
)

func isPrizePosition(x, y int) bool {
	for _, v := range prizePos {
		if v.X == x && v.Y == y {
			return true
		}
	}
	return false
}

func (gf *GameField) checkWinner(player Pair) (hasWon bool) {
	if isPrizePosition(player.X, player.Y) {
		return true
	}
	return false
}

func (gf *GameField) GetAvailableCells(playerIdx int) (cellsCoordinates []Pair) {
	var rowIdx int
	var secondPlayer *gfPlayer
	var player *gfPlayer
	cellsCoordinates = make([]Pair, 0)
	//TODO remove magic constants
	if playerIdx == 1 {
		player = &gf.p1
		rowIdx = 7
		secondPlayer = &gf.p2

	} else {
		player = &gf.p2
		rowIdx = 0
		secondPlayer = &gf.p1
	}

	//Get rows
	if player.pos == nil {
		for x := 0; x < fieldSize; x++ {
			if gf.field[rowIdx][x].isAvailable {
				cellsCoordinates = append(cellsCoordinates, Pair{x, rowIdx})
			}
		}
		return
	}
	//Set to left upper cell
	currCol := player.pos.X - 1
	currRow := player.pos.Y - 1

	for rowIdx := currRow; rowIdx < currRow+3; rowIdx++ {
		for colIdx := currCol; colIdx < currCol+3; colIdx++ {
			if rowIdx >= 0 && rowIdx < fieldSize && colIdx >= 0 && colIdx < fieldSize {
				if gf.field[rowIdx][colIdx].isAvailable {
					if secondPlayer.pos == nil {
						if (Pair{colIdx, rowIdx} != *player.pos) {
							cellsCoordinates = append(cellsCoordinates, Pair{colIdx, rowIdx})
						}
					} else {
						if (Pair{colIdx, rowIdx} != *player.pos) && (Pair{colIdx, rowIdx} != *secondPlayer.pos) {
							cellsCoordinates = append(cellsCoordinates, Pair{colIdx, rowIdx})
						}
					}

				}
			}
		}
	}

	return
}

//Поле для перемещения берется из регистров
func (gf *GameField) Move(playerIdx int) {
	var player *gfPlayer

	if playerIdx == 1 {
		player = &gf.p1
	} else {
		player = &gf.p2
	}

	//TODO этот метод должен получать ответ на regQuestion и проверять правильноть этого ответа
	if player.pos == nil {
		player.pos = &Pair{gf.regX, gf.regY}
	} else {
		player.pos.X = gf.regX
		player.pos.Y = gf.regY
	}

	//gf.Out <- event.Event{Etype: event.Move, Edata: player.id}
	return

}

func (gf *GameField) TryMovePlayer1(m message.Message) (e []event.Event, err error) {
	st := m.Payload.(map[string]interface{})
	nextX := int(st["x"].(float64))
	nextY := int(st["y"].(float64))
	if !gf.validateMoveCoordinates(&gf.p1, nextX, nextY) {
		err = errors.New(fmt.Sprintf("tried moving to invalid position x:%d ,y:%d", nextX, nextY))
		return
	}
	e, err = gf.tryMovePlayer(&gf.p1, nextX, nextY)
	if err != nil {
		logger.Error("TryMovePlayer1, tryMovePlayer returned error:", err)
		return
	}
	return
}

func (gf *GameField) TryMovePlayer2(m message.Message) (e []event.Event, err error) {
	st := m.Payload.(map[string]interface{})
	nextX := int(st["x"].(float64))
	nextY := int(st["y"].(float64))
	if !gf.validateMoveCoordinates(&gf.p2, nextX, nextY) {
		err = errors.New(fmt.Sprintf("tried moving to invalid position x:%d ,y:%d", nextX, nextY))
		return
	}
	e, err = gf.tryMovePlayer(&gf.p2, nextX, nextY)
	if err != nil {
		logger.Error("TryMovePlayer2, tryMovePlayer returned error:", err)
		return
	}
	return
}

func (gf *GameField) CheckIfMovesAvailable(playerId int) bool {

	availableCells := gf.GetAvailableCells(playerId)
	if len(availableCells) == 0 {
		return false
	} else {
		return true
	}

}

//Выполняет доставание вопроса из матрицы Игрового поля
func (gf *GameField) tryMovePlayer(player *gfPlayer, nextX int, nextY int) (e []event.Event, err error) {

	//destination := Pair{nextX, nextY}
	//TODO проверить, что destination isAvailable

	//Запись в регистр положения игрока, вопроса,
	gf.regY = nextY
	gf.regX = nextX

	//Здесь проверяем, если следущая клетка выигрышная

	if gf.checkWinner(Pair{nextX, nextY}) {
		e = make([]event.Event, 0)
		e = append(e, event.Event{Etype: event.WinPrize, Edata: nil})
		return
	}
	gf.regQuestion = *(gf.GetQuestionByCell(nextX, nextY))
	payload := struct {
		Question database.Question
		Time     int
	}{
		Question: gf.regQuestion,
		Time:     TurnDuration,
	}
	question, err := json.Marshal(payload)
	if err != nil {
		logger.Info("question unmarshal error")
		return
	}
	ms := string(question)

	e = make([]event.Event, 0)
	e = append(e, event.Event{Etype: event.Info, Edata: ms})
	return
}
