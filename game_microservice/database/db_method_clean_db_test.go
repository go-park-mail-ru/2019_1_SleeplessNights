package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"testing"
)

func TestCleanerDBForTests(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	pack := database.Pack{
		Theme:    "алгебра",
		IconPath: "math",
	}

	err = database.GetInstance().AddQuestionPack(pack)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	packs, err := database.GetInstance().GetPacksOfQuestions(1)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(packs) != 0 {
		t.Errorf("DB didn't cleaned up")
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}
