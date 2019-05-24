package database

func (db *dbManager) GetMessages(roomId uint64, since uint64, limit uint64) (messages []string, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT * FROM func_get_messages ($1::BIGINT, $2::BIGINT, $3::BIGINT)`,
		roomId, since, limit)
	if err != nil {
		logger.Errorf("Failed to get rows: %v", err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var str string
		err = rows.Scan(&str)
		if err != nil {
			logger.Errorf("Failed to get row: %v", err.Error())
			return
		}
		messages = append(messages, str)
	}

	err = rows.Err()
	if err != nil {
		logger.Errorf("Failed to scan: %v", err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Errorf("Failed to commit tx: %v", err.Error())
		return
	}
	return
}
