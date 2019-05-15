package database

func (db *dbManager) GetQuestions(packIDs []uint64) (questions []Question, qsFF []QuestionForFrontend, err error) {

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

		questions = append(questions, question)
		qsFF = append(qsFF, qFF)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
