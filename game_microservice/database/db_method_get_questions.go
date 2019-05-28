package database

func (db *dbManager) GetQuestions(packIDs []uint64) (questions []Question, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT * FROM func_get_questions($1)`, packIDs)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	var question Question
	var qFF QuestionForFrontend
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

		qFF.Text = question.Text
		qFF.Answers = question.Answers
		qFF.PackID = question.PackID
		jsonData, err := qFF.MarshalJSON()
		if err != nil {
			return nil, err
		}
		question.JSON = jsonData
		questions = append(questions, question)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
