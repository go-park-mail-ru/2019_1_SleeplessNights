package database

func (db *dbManager) GetUserSignature(email string) (id uint64, password, salt []byte, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT id, password, salt
	FROM public.func_get_user_by_email($1::CITEXT)`, email)
	err = row.Scan(
		&id,
		&password,
		&salt)
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
