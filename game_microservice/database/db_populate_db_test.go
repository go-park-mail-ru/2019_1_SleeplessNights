package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"testing"
)

func TestPopulateDatabase(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	database.GetInstance().PopulateDatabase()

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}
