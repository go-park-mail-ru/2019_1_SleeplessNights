package database

func (db *dbManager) AddQuestionPack(pack Pack) (err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM func_add_pack($1, $2)`, pack.Theme, pack.IconPath)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
