package database

import "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"

func (db *dbManager) GetUsers(page *services.PageData) (leaderBoardPage *services.LeaderBoardPage, err error) {
	tx, err := db.dataBase.Begin()
	if err != nil {
		logger.Error("Get users Method error, Begin transaction")

		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT id, email, nickname, avatar_path, rating, win_rate, matches
	FROM public.func_get_users($1::BIGINT, $2::BIGINT)`, page.Since, page.Limit)
	if err != nil {
		logger.Error("Get users Method error")
		return
	}
	defer rows.Close()

	var profiles []*services.Profile

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
			logger.Error("GetUsers db method Error")
			return
		}
		profiles = append(profiles, &profile)
	}
	leaderBoardPage = &services.LeaderBoardPage{
		Leaders: profiles,
	}
	err = rows.Err()
	if err != nil {
		logger.Error("GetUsers db method Error")
		return
	}
	logger.Info(profiles)
	err = tx.Commit()
	return
}
