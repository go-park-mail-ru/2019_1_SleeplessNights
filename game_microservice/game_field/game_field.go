package game_field

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/event"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/messge"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/questions"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"math"
)

var logger *log.Logger

func init () {
	logger = log.GetLogger("GameField")
}

const (
	fieldSize    = 8
	QuestionsNum = 60
)

var prizePos []pair

//В начале иры й игроков не существует никаких позиций, они находятся как бы вне поля

func init() {
	prizePos = []pair{{3, 3}, {3, 4}, {4, 3}, {4, 4}}
}

type gameCell struct {
	isAvailable bool
	question    *questions.Question
}

type pair struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type gfPlayer struct {
	pos *pair //Поставил указатель на pair, чтобы pos поддерживала значение nil (начальные условия)
	id  uint64
}

type GameField struct {
	field [fieldSize][fieldSize]gameCell
	p1    gfPlayer
	p2    gfPlayer
	//Out   []event.Event

	regX        int
	regY        int
	regQuestion questions.Question

	//Тут уровень абстракции уже достаточно постой для понимания, поэтому оставляю реализацию на ваше усмотрение
	//По ответственности, если навскидку, игровое поле должно:
	//* Знать расположение всех вопросов
	//* Знать координаты всех игроков (при этом самих игроков знать не должно)
	//* Знать состояние каждой клетки поля: доступна/недоступна для хода
	//* Уметь проверять, что кто-то из игроков уже достиг центра и победил
	//* Уметь проверять, что кто-то из игроков оказался недостижим от центра поля и проиграл (задачка со звёздочкой на алгоритм Дейкстры)
	//TODO develop
}

func isPrizePosition(x, y int) bool {
	for _, v := range prizePos {
		if v.X == x && v.Y == y {
			return true
		}
	}
	return false
}

func (gf *GameField) Build(qArray [QuestionsNum]questions.Question) {
	qSlice := qArray[:]
	index := 0
	for rowIdx, row := range gf.field {
		for colIdx := range row {
			if isPrizePosition(rowIdx, colIdx) {
				gf.field[rowIdx][colIdx] = gameCell{true, nil}
			} else {
				gf.field[rowIdx][colIdx] = gameCell{true, &qSlice[index]}
				index++
			}
		}
	}
	gf.p1.pos = nil
	gf.p2.pos = nil
}

func (gf *GameField) checkWinner(player pair) (hasWon bool) {
	if isPrizePosition(player.X, player.Y) {
		return true
	}
	return false
}

//Пока не трогать
func (gf *GameField) checkRouteAvailable(player pair) (isAvailable bool) {
	//Call some Dijkstra algorithm
	//Для заглушки игрок всегда достижим относительно приза
	return true
}

/*

 x0 x1 x2 x3 x4 x5 x6 x7
 __ __ __ __ __ __ __ __
|__|__|__|__|P2|__|__|__|  y0
|__|__|__|__|_X|_X|__|__|  y1
|__|__|__|__|_X|_X|_X|_X|  y2
|__|__|__|Pr|Pr|p1|__|__|  y3
|__|__|__|Pr|Pr|__|__|__|  y4
|__|__|__|__|_X|__|__|__|  y5
|__|__|__|__|_X|__|__|__|  y6
|__|__|__|P1|__|__|__|__|  y7

*/

func (gf *GameField) GetAvailableCells(playerIdx int) (cellsCoordinates []pair) {
	var rowIdx int
	var secondPlayer *gfPlayer
	var player *gfPlayer

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
			cellsCoordinates = append(cellsCoordinates, pair{x, rowIdx})
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
					if (pair{colIdx, rowIdx} != *player.pos) && (pair{colIdx, rowIdx} != *secondPlayer.pos) {
						cellsCoordinates = append(cellsCoordinates, pair{colIdx, rowIdx})
					}
				}
			}
		}
	}

	return
}

func (gf *GameField) Move(player *gfPlayer) {
	//TODO этот метод должен получать ответ на regQuestion и проверять правильноть этого ответа
	player.pos.X = gf.regX
	player.pos.Y = gf.regY

	if gf.checkWinner(*player.pos) {
		//gf.Out <- event.Event{Etype: event.WinPrize, Edata: player.id}
		return
	}

	//gf.Out <- event.Event{Etype: event.Move, Edata: player.id}
	return

}

func (gf *GameField) TryMovePlayer1(m messge.Message) (e []event.Event, err error) {
	nextX := m.Payload.(*messge.Coordinates).X
	nextY := m.Payload.(*messge.Coordinates).Y

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

func (gf *GameField) TryMovePlayer2(m messge.Message) (e []event.Event, err error) {
	nextX := m.Payload.(*messge.Coordinates).X
	nextY := m.Payload.(*messge.Coordinates).Y

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

//Выполняет доставание вопроса из матрицы Игрового поля
func (gf *GameField) tryMovePlayer(player *gfPlayer, nextX int, nextY int) (e []event.Event, err error) {

	//destination := pair{nextX, nextY}
	//TODO проверить, что destination isAvailable

	//Запись в регистр положения игрока, вопроса,
	gf.regY = nextY
	gf.regX = nextX

	//Пока не трогать
	if !gf.checkRouteAvailable(*gf.p1.pos) {
		//TODO отправить Event Loose для текущего игрока и Event Win для второго игрока

		//TODO переместить в начало метода GetAvailableCells
	}

	ms := struct {
		question string
	}{
		gf.GetQuestionByCell(nextX, nextY).QuestionJson,
	}

	e = make([]event.Event, 0)
	e = append(e, event.Event{Etype: event.Move, Edata: ms})
	return
}

func (gf *GameField) GetQuestionByCell(x, y int) (question questions.Question) {
	logger.Infof("GetQuestionByCell x:%d,y:%d ", x, y)
	question = *(gf.field[y][x].question)
	return
}

func (gf *GameField) validateMoveCoordinates(player *gfPlayer, nextX int, nextY int) bool {
	nextPos := pair{nextX, nextY}
	//Убрать Валидацию поля в GameField

	if nextX > fieldSize || nextY > fieldSize || nextX < 0 || nextY < 0 {
		logger.Error("Coordinate validator, error:invalid next coordinates")
		return false
	}

	if player.pos == nil {
		return true
	}

	if math.Abs(float64(player.pos.X-nextX)) > 1 || math.Abs(float64(player.pos.Y-nextY)) > 1 {
		logger.Error("Coordinate validator, error: player trie to reach cell", nextPos)
		return false
	}

	if gf.p1.pos == nil || gf.p2.pos == nil {
		return true
	}

	if (*gf.p1.pos) == nextPos || (*gf.p2.pos) == nextPos {
		logger.Errorf("Desired Position is another's players position p1:%v , p2:%v , desiredPos:%v", gf.p1.pos, gf.p2.pos, nextPos)
		return false
	}

	return true
}

func (gf *GameField) CheckAnswer(answerIdx int) bool {
	if !gf.validateAnswerId(answerIdx) {
		return false
	}
	if gf.regQuestion.CorrectAnswerId == answerIdx {
		return true
	}
	return false
}

func (gf *GameField) validateAnswerId(answerId int) bool {
	//Убрать Валидацию поля в GameField
	if answerId > 3 || answerId < 0 {
		logger.Error("validateAnswerId, error:AnswerId")
		return false
	}
	return true
}
