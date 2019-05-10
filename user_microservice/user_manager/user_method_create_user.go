package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"golang.org/x/net/context"
)

const defaultAvatar = "default_avatar.jpg"

func (us *userManager) CreateUser(ctx context.Context, in *services.NewUserData) (*services.User, error) {

	logger.Debug("Enter_OK")

	salt, err := MakeSalt()
	if err != nil {
		return nil, err
	}
	logger.Debug("MakeSalt_OK")

	password := MakePasswordHash(in.Password, salt)

	logger.Debug("MakePasswordHash_OK")

	user, err := database.GetInstance().AddUser(in.Email, in.Nickname, defaultAvatar, password, salt)
	if err != nil {
		return nil, err
	}

	logger.Debug("AddUser_OK")

	return &user, nil
}
