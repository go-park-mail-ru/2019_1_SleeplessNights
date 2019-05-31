package database

func (db *dbManager) AddMessage(talkerId uint64, roomId uint64, payload []byte) (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM func_add_message ($1::BIGINT, $2::BIGINT, $3::JSON)`,
		talkerId, roomId, payload)
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
