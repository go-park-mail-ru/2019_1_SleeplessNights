package database

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
)

func (db *dbManager) GetProfile(userID uint64) (profile services.Profile, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT id, email, nickname, avatar_path, rating, win_rate, matches
	FROM public.func_get_user_by_id($1::BIGINT)`, userID)
	var user services.User
	profile.User = &user
	err = row.Scan(
		&profile.User.Id,
		&profile.User.Email,
		&profile.User.Nickname,
		&profile.User.AvatarPath,
		&profile.Rating,
		&profile.WinRate,
		&profile.Matches)
	if err != nil {
		logger.Errorf("Failed to get row: %v", err.Error())
		return
	}

	err = tx.Commit()
	if err !=  nil {
		logger.Errorf("Failed to commit tx: %v", err.Error())
		return
	}
	return
}
