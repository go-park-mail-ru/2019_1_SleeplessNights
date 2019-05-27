package database

import "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"

func (db *dbManager) AddUser(email, nickname, avatarPath string, password, salt []byte) (user services.User, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT id, email, nickname, avatar_path
	FROM public.func_add_user($1::CITEXT, $2::BYTEA, $3::BYTEA, $4::TEXT, $5::TEXT)`,
		email, password, salt, nickname, avatarPath)
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
