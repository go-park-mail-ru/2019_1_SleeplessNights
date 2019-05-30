package game_field_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/game_field"
	"testing"
)

func TestBuild(t *testing.T) {
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

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}
