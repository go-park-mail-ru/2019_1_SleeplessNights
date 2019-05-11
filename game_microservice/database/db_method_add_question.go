package database

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database/models"
)

func (db *dbManager) AddQuestion(question models.Question) (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM func_add_question($1, $2, $3, $4)`,
		question.Answers, question.Correct, question.Text, question.PackID)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
