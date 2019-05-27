package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/jackc/pgx"
	"golang.org/x/net/context"
)

func (us *userManager) GetProfile(ctx context.Context, in *services.User) (*services.Profile, error) {
	profile, err := database.GetInstance().GetProfile(in.Id)
	if _err, ok := err.(pgx.PgError); ok {
		logger.Errorf("Failed to get profile: %v", err.Error())
		err = handlerError(_err)
		return nil, err
	}
	return &profile, nil
}

func (us *userManager) UpdateProfile(ctx context.Context, in *services.User) (*services.User, error) {
	err := database.GetInstance().UpdateUser(in)
	if _err, ok := err.(pgx.PgError); ok {
		logger.Errorf("Failed to update user: %v", err.Error())
		err = handlerError(_err)
		return nil, err
	}
	return in, nil
}
