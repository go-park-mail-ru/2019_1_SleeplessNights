package database

func (db *dbManager) GetPacksOfQuestions(number int) (packs []Pack, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT * FROM func_get_packs($1)`, number)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	var pack Pack
	for rows.Next() {

		err = rows.Scan(
			&pack.ID,
			&pack.IconPath,
			&pack.Theme)
		if err != nil {
			return
		}
		packs = append(packs, pack)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
