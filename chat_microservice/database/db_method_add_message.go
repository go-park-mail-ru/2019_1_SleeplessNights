package database

func (db *dbManager) AddMessage(userId uint64, roomId uint64, payload []byte) (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM func_add_message ($1, $2, $3)`,
		userId, roomId, payload)
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
