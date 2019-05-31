package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/jackc/pgx"
	"golang.org/x/net/context"
)

func (us *userManager) GetLeaderBoardPage(ctx context.Context, in *services.PageData) (*services.LeaderBoardPage, error) {
	page, err := database.GetInstance().GetUsers(in)
	if _err, ok := err.(pgx.PgError); ok {
		logger.Errorf("Failed to get users: %v", err.Error())
		err = handlerError(_err)
		return nil, err
	}
	return page, nil
}
