package game_field_test

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/game_field"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/message"
	"testing"
)

func TestGetAvailableCells (t *testing.T){
	gf := game_field.GameField{}
	gf.GetAvailableCells(1)
	gf.SetPlayerFirstPosition()
	gf.GetAvailableCells(1)
	gf.GetAvailableCells(2)
}

func TestMOVE (t *testing.T){
	gf := game_field.GameField{}
	gf.Move(1)
	gf.Move(2)
}

func TestTryMovePlayer1 (t *testing.T){
	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	database.GetInstance().PopulateDatabase()

	gf := game_field.GameField{}

	question, err := database.GetInstance().GetQuestions([]uint64{1,2,3,4,5,6})

	if err != nil {
		t.Errorf(err.Error())
	}

	gf.Build(question)

	result, _ := json.Marshal(message.Message{
		Title:   message.GoTo,
		Payload: message.Coordinates{X: 0, Y: 1},
	})
	var GoToMessage message.Message
	_ = json.Unmarshal(result, &GoToMessage)
	_, err = gf.TryMovePlayer1(GoToMessage)
	if err != nil {
		t.Error(err.Error())
	}

	result, _ = json.Marshal(message.Message{
		Title:   message.GoTo,
		Payload: message.Coordinates{X: -1, Y: -1},
	})
	_ = json.Unmarshal(result, &GoToMessage)
	_, err = gf.TryMovePlayer1(GoToMessage)
	if err == nil {
		t.Errorf("Didn't retunr error!!!")
	}

	result, _ = json.Marshal(message.Message{
		Title:   message.GoTo,
		Payload: message.Coordinates{X: 5, Y: 6},
	})
	_ = json.Unmarshal(result, &GoToMessage)
	_, err = gf.TryMovePlayer1(GoToMessage)
	if err != nil {
		t.Error(err.Error())
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestTryMovePlayer2 (t *testing.T){
	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	database.GetInstance().PopulateDatabase()

	gf := game_field.GameField{}

	question, err := database.GetInstance().GetQuestions([]uint64{1,2,3,4,5,6})

	if err != nil {
		t.Errorf(err.Error())
	}

	gf.Build(question)

	result, _ := json.Marshal(message.Message{
		Title:   message.GoTo,
		Payload: message.Coordinates{X: 0, Y: 1},
	})
	var GoToMessage message.Message
	_ = json.Unmarshal(result, &GoToMessage)
	_, err = gf.TryMovePlayer2(GoToMessage)
	if err != nil {
		t.Error(err.Error())
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}