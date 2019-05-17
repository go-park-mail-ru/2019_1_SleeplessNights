package database_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/jackc/pgx"
	"testing"
)

func TestAddQuestionSuccessful(t *testing.T) {

	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}

	pack := database.Pack{
		Theme:    "алгебра",
		IconPath: "math_icon.jpg",
	}

	err = database.GetInstance().AddQuestionPack(pack)
	if err != nil {
		t.Errorf(err.Error())
	}

	question := database.Question{
		Answers: []string{"sdf", "sdf"},
		Correct: 1,
		PackID:  1,
		Text:    "sdfsdf",
	}

	err = database.GetInstance().AddQuestion(question)
	if err != nil {
		t.Errorf(err.Error())
	}

	ids := []uint64{1}

	questions, err := database.GetInstance().GetQuestions(ids)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(questions) != 1 {
		t.Errorf("DB returned wrong len of questions:\ngot %v\nwant %v", len(questions), 1)
	} else if questions[0].Text != question.Text {
		t.Errorf("DB returned wrong questions text:\ngot %v\nwant %v",
			questions[0].Text, question.Text)
	}
}

func TestAddQuestionUnsuccessful(t *testing.T) {

	question := database.Question{
		Answers: []string{"sdf", "sdf"},
		Correct: 1,
		PackID:  2,
		Text:    "sdfsdf",
	}

	expected := "23505"

	err := database.GetInstance().AddQuestion(question)
	if err == nil {
		t.Errorf("DB didn't return error!")
	} else if err, ok := err.(pgx.PgError); ok && err.Code == expected {
		t.Errorf("DB returned wrong error:\ngot %v\nwant %v",
			err.Code, expected)
	}

	err = database.GetInstance().CleanerDBForTests()
	if err != nil {
		t.Errorf(err.Error())
	}
}
