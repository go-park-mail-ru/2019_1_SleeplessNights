package database

func (db *dbManager) UpdateUser(uid uint64, nickname string, avatarPath string) (id uint64, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT * FROM func_update_user ($1::BIGINT, $2::TEXT, $3::TEXT)`,
		uid, nickname, avatarPath)
	err = row.Scan(&id)
	if err != nil {
		logger.Errorf("Failed to get row: %v", err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Errorf("Failed to commit tx: %v", err.Error())
		return
	}
	return
}
