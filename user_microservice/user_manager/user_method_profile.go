package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"golang.org/x/net/context"
)

func (us *userManager)GetProfile(ctx context.Context, in *services.User)(*services.Profile, error) {
	profile, err := database.GetInstance().GetProfile(in.Id)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (us *userManager)UpdateProfile(ctx context.Context, in *services.User)(*services.User, error) {
	err := database.GetInstance().UpdateUser(in)
	if err != nil {
		return nil, err
	}
	return in, nil
}