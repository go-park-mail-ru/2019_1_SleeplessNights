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
	playersPos = []pair{{4, 0}, {3, 7}}
}

type gameCell struct {
	isAvailable bool
	question    questions.Question
}

type pair struct {
	x int
	y int
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
		if v.x == x && v.y == y {
			return true
		}
	}
	return false
}

func isPlayersPosition(x, y int) bool {
	for _, v := range playersPos {
		if v.x == x && v.y == y {
			return true
		}
	}
	return false
}

func (gf *GameField) Build(qSlice []questions.Question) {
	for rowIdx, row := range gf.field {
		for colIdx := range row {
			if isPlayersPosition(rowIdx, colIdx) || isPrizePosition(rowIdx, colIdx) {
				gf.field[rowIdx][colIdx] = gameCell{true, nil}

			} else {
				rand.Seed(time.Now().UnixNano())
				index := rand.Intn(questionsNum)
				gf.field[rowIdx][colIdx] = gameCell{true, qSlice[index]}
				qSlice = append(qSlice[:index], qSlice[index+1:]...)
			}
		}
	}
}

func (gf *GameField) checkWinner(player pair) (hasWon bool) {
	if isPrizePosition(player.x, player.y) {
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

func (gf *GameField) MovePlayer(playerId uint64, nextX int, nextY int, ch chan event.Event) {

	var playerToMove *gfPlayer
	var playerOther *gfPlayer

	if gf.p1.id == playerId {
		playerToMove = &gf.p1
		playerOther = &gf.p2

	} else if gf.p2.id != playerId {
		playerToMove = &gf.p2
		playerOther = &gf.p1

	} else {
		ch <- event.Event{Etype: event.Error, Edata: nil}
		return
	}

	if nextX > 7 || nextY > 7 || nextX < 0 || nextY < 0 {
		//If Wanna go outside the board
		ch <- event.Event{Etype: event.Warning, Edata: playerToMove.id}
		return
	}

	if (pair{nextX, nextY}) == playerToMove.pos {
		// Tell about Same position or skip his turn?
		ch <- event.Event{Etype: event.Warning, Edata: playerToMove.id}
		return
	}

	if gf.checkWinner(playerToMove.pos) {
		ch <- event.Event{Etype: event.WinPrize, Edata: playerToMove.id}
		ch <- event.Event{Etype: event.Lose, Edata: playerOther.id}
		return
	}

	//Пока не трогать
	if gf.checkRouteAvailable(gf.p1.pos) {
	}

	//If all things are as usual then Move player and Send ? question ? and player ID

	ch <- event.Event{Etype: event.Move, Edata: {playerToMove.id}}
	return
}
