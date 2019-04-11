package game_field

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/event"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/questions"
	"math/rand"
	"time"
)

const (
	fieldSize    = 8
	questionsNum = 60
)

var prizePos []pair
var playersPos []pair

func init() {
	prizePos = []pair{{3, 3}, {3, 4}, {4, 3}, {4, 4}}
}

type gameCell struct {
	isAvailable bool
	question    questions.Question
}

type pair struct {
	X int `json:"X"`
	Y int `json:"Y"`
}

type gfPlayer struct {
	pos pair
	id  uint64
}

type GameField struct {
	field [fieldSize][fieldSize]gameCell
	p1    gfPlayer
	p2    gfPlayer
	Out   chan event.Event

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

func (gf *GameField) Build(qArray [questionsNum]questions.Question) {
	qSlice := qArray[:]
	for rowIdx, row := range gf.field {
		for colIdx := range row {
			if isPrizePosition(rowIdx, colIdx) {
				gf.field[rowIdx][colIdx] = gameCell{true, nil}
			} else {
				rand.Seed(time.Now().UnixNano())
				index := rand.Intn(questionsNum)
				gf.field[rowIdx][colIdx] = gameCell{true, qSlice[index]}
				qSlice = append(qSlice[:index], qSlice[index+1:]...)
			}
		}
	}
	gf.p1.pos = pair{-1, -1}

	gf.p2.pos = pair{-1, -1}
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
	return
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
func (gf *GameField) TryMovePlayer1(nextX int, nextY int) {
	gf.tryMovePlayer(&gf.p1, nextX, nextY)
}

func (gf *GameField) TryMovePlayer2(nextX int, nextY int) {
	gf.tryMovePlayer(&gf.p2, nextX, nextY)
}

func (gf *GameField) GetAvailableCells(player *gfPlayer) (cellsCoordinates []pair) {
	var rowIdx int
	var secondPlayer *gfPlayer

	if player.id == gf.p1.id {
		rowIdx = 7
		secondPlayer = &gf.p2

	} else {
		rowIdx = 0
		secondPlayer = &gf.p1
	}

	//Get rows
	if (player.pos == pair{-1, -1}) {
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
			if rowIdx >= 0 && rowIdx <= 7 && colIdx >= 0 || colIdx <= 7 {
				if gf.field[rowIdx][colIdx].isAvailable {
					if (pair{colIdx, rowIdx} != player.pos) && (pair{colIdx, rowIdx} != secondPlayer.pos) {
						cellsCoordinates = append(cellsCoordinates, pair{colIdx, rowIdx})
					}
				}
			}
		}
	}

	return
}

func (gf *GameField) Move(player *gfPlayer) {

	player.pos.X = gf.regX
	player.pos.Y = gf.regY

	if gf.checkWinner(player.pos) {
		gf.Out <- event.Event{Etype: event.WinPrize, Edata: player.id}
		return
	}

	gf.Out <- event.Event{Etype: event.Move, Edata: player.id}
	return

}

func (gf *GameField) tryMovePlayer(player *gfPlayer, nextX int, nextY int) {

	if nextX > 7 || nextY > 7 || nextX < 0 || nextY < 0 {
		//Make logging of invalid commands instead of sending  events
		gf.Out <- event.Event{Etype: event.Warning, Edata: player.id}
		return
	}

	if (pair{nextX, nextY}) == player.pos {
		//Make logging of invalid commands instead of sending  events
		gf.Out <- event.Event{Etype: event.Warning, Edata: player.id}
		return
	}

	gf.regY = nextY
	gf.regX = nextX

	//Пока не трогать
	if gf.checkRouteAvailable(gf.p1.pos) {

	}

	ms := struct {
		playerId uint64
		question questions.Question
	}{
		player.id, gf.field[nextY][nextX].question,
	}
	gf.Out <- event.Event{Etype: event.Move, Edata: ms}
	return
}

func (gf *GameField) GetQuestionByCell(x, y int) (question string) {
	question = gf.field[y][x].question.QuestionJson
	return
}
