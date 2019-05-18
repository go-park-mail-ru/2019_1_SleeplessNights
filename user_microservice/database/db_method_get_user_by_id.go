package database

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
)

func (db *dbManager) GetUserByID(userID uint64) (user services.User, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT id, email, nickname, avatar_path
	FROM public.func_get_user_by_id($1::BIGINT)`, userID)
	err = row.Scan(
		&user.Id,
		&user.Email,
		&user.Nickname,
		&user.AvatarPath)
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
