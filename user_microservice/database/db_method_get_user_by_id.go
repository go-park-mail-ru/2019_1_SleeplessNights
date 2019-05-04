package database

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
)

func (db *dbManager) GetUserByID(userID uint64) (user services.User, err error) {

	tx, err := db.dataBase.Begin()
	if err != nil {
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
		return
	}

	err = tx.Commit()
	return
}
