package user_manager

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/jackc/pgx"
)

func (us *userManager) UpdateStats(ctx context.Context, in *services.MatchResults) (*services.Nothing, error){
	err := database.GetInstance().UpdateStats(in)
	if _err, ok := err.(pgx.PgError); ok {
		logger.Errorf("Failed to update stats: %v", err.Error())
		err = handlerError(_err)
		return nil, err
	}
	var n *services.Nothing
	return n, nil
}
