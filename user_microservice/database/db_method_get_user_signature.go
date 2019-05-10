package database

func (db *dbManager) GetUserSignature(email string) (id uint64, password, salt []byte, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT id, password, salt
	FROM public.func_get_user_by_email($1::CITEXT)`, email)
	err = row.Scan(&id, &password, &salt)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
