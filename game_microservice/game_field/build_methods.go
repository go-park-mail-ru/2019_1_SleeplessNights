package game_field

import "github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"

func (gf *GameField) Build(qArray []database.Question) {
	Shuffle(&qArray)
	qSlice := qArray

	index := 0
	for rowIdx, row := range gf.field {
		for colIdx := range row {
			if isPrizePosition(rowIdx, colIdx) {
				gf.field[rowIdx][colIdx] = gameCell{true, 0, nil}
			} else {
				gf.field[rowIdx][colIdx] = gameCell{true, 0, &qSlice[index]}
				index++
			}
		}
	}
	gf.p1.pos = nil
	gf.p2.pos = nil
}
