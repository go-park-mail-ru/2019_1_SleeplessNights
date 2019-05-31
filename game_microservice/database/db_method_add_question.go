package database

func (db *dbManager) AddQuestion(question Question) (err error) {

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
