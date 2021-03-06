package database

import "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"


func (db *dbManager) GetUsers(limit uint64) (profiles []*services.Profile, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Errorf("Failed to begin tx: %v", err.Error())
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT id, email, nickname, avatar_path, rating, win_rate, matches
	FROM public.func_get_users($1::BIGINT)`, limit)
	if err != nil {
		logger.Errorf("Failed to get rows: %v", err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var profile services.Profile
		var user services.User
		profile.User = &user
		err = rows.Scan(
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
		profiles = append(profiles, &profile)
	}

	err = rows.Err()
	if err != nil {
		logger.Errorf("Failed to scan: %v", err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Errorf("Failed to commit tx: %v", err.Error())
		return
	}
	return
}
