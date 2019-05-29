package game_field_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/game_field"
	"testing"
)

func TestResetPlayersPositions(t *testing.T) {
	gf := game_field.GameField{}
	gf.ResetPlayersPositions()
}

func TestGetCurrentState(t *testing.T) {
	gf := game_field.GameField{}
	gf.GetCurrentState()
}