package database

func (db *dbManager) CleanerDBForTests() (err error) {
	//TODO remove?
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM public.func_clean_db()`)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
