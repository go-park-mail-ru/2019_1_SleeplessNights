package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"golang.org/x/net/context"
)

func (auth *authManager)GetLeaderBoardPage(ctx context.Context, in *services.PageData)(*services.LeaderBoardPage, error) {
	users, err := database.GetInstance().GetUsers(in)
	if err != nil {
		return nil, err
	}

	page := services.LeaderBoardPage{
		Users: users,
	}

	return &page, nil
}

