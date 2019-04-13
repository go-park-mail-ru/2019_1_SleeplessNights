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
//В начале иры й игроков не существует никаких позиций, они находятся как бы вне поля

func init() {
	prizePos = []pair{{3, 3}, {3, 4}, {4, 3}, {4, 4}}
}

type gameCell struct {
	isAvailable bool
	question    questions.Question
}

type pair struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type gfPlayer struct {
	pos *pair//Поставил указатель на pair, чтобы pos поддерживала значение nil (начальные условия)
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
				gf.field[rowIdx][colIdx] = gameCell{true,nil}
			} else {
				rand.Seed(time.Now().UnixNano())//TODO еспли я правильно помню, то seed нужно скормить 1 раз, а не в цикле. Погугли, пожалуйста
				index := rand.Intn(len(qSlice))
				gf.field[rowIdx][colIdx] = gameCell{true, qSlice[index]}
				qSlice = append(qSlice[:index], qSlice[index+1:]...)
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
		gf.Out <- event.Event{Etype: event.WinPrize, Edata: player.id}
		return
	}

	gf.Out <- event.Event{Etype: event.Move, Edata: player.id}
	return

}

func (gf *GameField) tryMovePlayer(player *gfPlayer, nextX int, nextY int) {

	if nextX >= fieldSize || nextY >= fieldSize || nextX < 0 || nextY < 0 {
		//TODO Make logging of invalid commands instead of sending  events
		gf.Out <- event.Event{Etype: event.Warning, Edata: player.id}
		return
	}

	destination := pair{nextX, nextY}
	//TODO проверить, что destination isAvailable
	//TODO проверить, что модуль разницы по обеим координатом между player.pos и destination не превышает еденицы
	if destination == *player.pos {
		//TODO Make logging of invalid commands instead of sending  events
		gf.Out <- event.Event{Etype: event.Warning, Edata: player.id}
		return
	}

	gf.regY = nextY
	gf.regX = nextX

	//Пока не трогать
	if !gf.checkRouteAvailable(*gf.p1.pos) {
		//TODO отправить Event Loose для текущего игрока и Event Win для второго игрока
		//TODO переместить в начало метода GetAvailableCells
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
	//А почему y, x? Разве не gf.field[x][y]?
	question = gf.field[y][x].question.QuestionJson
	return
}
