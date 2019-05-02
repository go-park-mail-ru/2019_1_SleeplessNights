package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"golang.org/x/net/context"
)

const defaultAvatar = "default_avatar.jpg"

func (auth *authManager)CreateUser(ctx context.Context, in *services.NewUserData)(*services.SessionToken, error) {
	salt, err := MakeSalt()
	if err != nil {
		return nil, err
	}

	password := MakePasswordHash(in.Password, salt)

	err = database.GetInstance().AddUser(in.Email, in.Nickname, defaultAvatar, password, salt)
	if err != nil {
		return nil, err
	}

	signature := services.UserSignature{
		Email:    in.Email,
		Password: in.Password,
	}
	token, err := auth.MakeToken(ctx, &signature)
	if err != nil {
		return nil, err
	}

	return token, nil
}