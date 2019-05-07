package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"golang.org/x/net/context"
)

func (us *userManager) GetUserById(ctx context.Context, in *services.UserId) (*services.User, error) {
	user, err := database.GetInstance().GetUserByID(in.Id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
