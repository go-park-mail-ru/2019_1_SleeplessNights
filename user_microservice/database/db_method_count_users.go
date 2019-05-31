package database

func (db *dbManager) GetLenUsers() (len int, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT COUNT(*) FROM public.users`)
	err = row.Scan(&len)
	if err != nil {
		logger.Errorf("Failed to get row: %v", err.Error())
		return
	}

	err = tx.Commit()
	if err !=  nil {
		logger.Errorf("Failed to commit tx: %v", err.Error())
		return
	}
	return
}