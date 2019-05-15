package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"testing"
)

func TestAddQuestionPack(t *testing.T) {

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

	packs, err := database.GetInstance().GetPacksOfQuestions(1);
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(packs) != 1 {
		t.Errorf("DB returned wrong len of packs:\ngot %v\nwant %v", len(packs), 1)
	}

	if packs[0].Theme != pack.Theme {
		t.Errorf("DB returned wrong theme:\ngot %v\nwant %v",
			packs[0].Theme, pack.Theme)
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}
