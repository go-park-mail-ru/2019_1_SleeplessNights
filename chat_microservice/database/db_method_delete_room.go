package database

func (db *dbManager) DeleteRoom(roomId uint64) (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM func_delete_room($1::BIGINT)`, roomId)
	if err != nil {
		logger.Errorf("Failed to exec: %v", err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Errorf("Failed to commit tx: %v", err.Error())
		return
	}
	return
}
