package database

import "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"

func (db *dbManager) UpdateUser(user *services.User) (err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`SELECT * FROM public.func_update_user($1::TEXT, $2::TEXT, $3::BIGINT)`,
		user.Nickname, user.AvatarPath, user.Id)
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}
