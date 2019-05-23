package database

func (db *dbManager) GetRoomsIds() (roomIds []uint64, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT * FROM func_get_rooms()`)
	if err != nil {
		logger.Errorf("Failed to get rows: %v", err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var roomId uint64
		err = rows.Scan(&roomId)
		if err != nil {
			logger.Errorf("Failed to get row: %v", err.Error())
			return
		}
		roomIds = append(roomIds, roomId)
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