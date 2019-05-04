package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"golang.org/x/net/context"
)

const defaultAvatar = "default_avatar.jpg"

func (us *userManager) CreateUser(ctx context.Context, in *services.NewUserData) (*services.User, error) {
	salt, err := MakeSalt()
	if err != nil {
		return nil, err
	}

	password := MakePasswordHash(in.Password, salt)

	user, err := database.GetInstance().AddUser(in.Email, in.Nickname, defaultAvatar, password, salt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
