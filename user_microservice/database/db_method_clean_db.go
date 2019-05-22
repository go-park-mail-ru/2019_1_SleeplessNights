package database

func (db *dbManager) CleanerDBForTests() (err error) {
	//TODO remove?
	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM public.func_clean_user_db()`)
	if err != nil {
		logger.Errorf("Failed to exec: %v", err.Error())
		return
	}

	err = tx.Commit()
	if err !=  nil {
		logger.Errorf("Failed to commit tx: %v", err.Error())
		return
	}
	return
}
