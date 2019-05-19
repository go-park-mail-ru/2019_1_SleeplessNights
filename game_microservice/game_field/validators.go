package game_field

import "math"

func (gf *GameField) validateMoveCoordinates(player *gfPlayer, nextX int, nextY int) bool {
	nextPos := Pair{nextX, nextY}
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

func (gf *GameField) validateAnswerId(answerId int) bool {
	//Убрать Валидацию поля в GameField
	if answerId > 3 || answerId < 0 {
		logger.Error("validateAnswerId, error:AnswerId")
		return false
	}
	return true
}
