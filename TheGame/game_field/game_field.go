package game_field

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/event"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/questions"
)

const (
	fieldSize = 8
)

type gameCell struct {
	isAvailable bool
	question    *questions.Question
	x           int
	y           int
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
	field    [fieldSize][fieldSize]gameCell
	p1       gfPlayer
	p2       gfPlayer
	prizePos pair
	Out      chan event.Event

	//Тут уровень абстракции уже достаточно постой для понимания, поэтому оставляю реализацию на ваше усмотрение
	//По ответственности, если навскидку, игровое поле должно:
	//* Знать расположение всех вопросов
	//* Знать координаты всех игроков (при этом самих игроков знать не должно)
	//* Знать состояние каждой клетки поля: доступна/недоступна для хода
	//* Уметь проверять, что кто-то из игроков уже достиг центра и победил
	//* Уметь проверять, что кто-то из игроков оказался недостижим от центра поля и проиграл (задачка со звёздочкой на алгоритм Дейкстры)
	//TODO develop
}

func (gf *GameField) checkWinner(player pair) (isWon bool) {
	if player == gf.prizePos {
		return true
	}
	return false
}

func (gf *GameField) checkRouteAvailable(player pair) (isAvailable bool) {

	return
}

/*

 x0 x1 x2 x3 x4 x5 x6 x7
 __ __ __ __ __ __ __ __
|p2|__|__|__|__|__|__|__|  y0
|__|__|__|__|_X|_X|__|__|  y1
|__|__|__|__|_X|_X|_X|_X|  y2
|__|__|__|__|_X|p1|__|__|  y3
|__|__|__|Pr|_X|__|__|__|  y4
|__|__|__|__|_X|__|__|__|  y5
|__|__|__|__|_X|__|__|__|  y6
|__|__|__|__|_X|__|__|__|  y7

*/

func (gf *GameField) MovePlayer1(x int, y int, ch chan event.Event) {
	nextPos := pair{x, y}
	if x > 7 || y > 7 {
		ch <- event.Event{}
		return
	}
	if nextPos == gf.p1.pos {
		ch <- event.Event{}
		return
	}
	if gf.checkRouteAvailable(gf.p1.pos) {

	}

	return
}

func (gf *GameField) MovePlayer2(x int, y int, ch chan event.Event) {
	return
}
