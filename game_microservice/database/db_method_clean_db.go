package database

func (db *dbManager) CleanerDBForTests() (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM public.func_clean_game_db()`)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
