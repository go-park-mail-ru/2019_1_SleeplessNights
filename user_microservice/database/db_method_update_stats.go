package database

import "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"

func (db *dbManager) UpdateStats(in *services.MatchResults) (err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM public.func_update_stats($1::BIGINT, $2::BIGINT, $3::BIGINT)`,
		in.Winner, in.Loser, in.WinnerRating)
	if err != nil {
		logger.Errorf("Failed to exec: %v", err.Error())
		return
	}

	err = tx.Commit()
	if err !=  nil {
		logger.Errorf("Failed to commit tx: %v", err.Error())
		return
	}
	return
}
