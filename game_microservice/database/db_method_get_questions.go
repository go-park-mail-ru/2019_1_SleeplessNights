package database

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database/models"
)

func (db *dbManager) GetQuestions(ids []int) (questions []models.Question, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT * FROM func_get_questions($1)`, ids)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	var question models.Question
	for rows.Next() {
		err = rows.Scan(
			&question.ID,
			&question.Answers,
			&question.Correct,
			&question.Text,
			&question.PackID)
		if err != nil {
			return
		}
		questions = append(questions, question)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
