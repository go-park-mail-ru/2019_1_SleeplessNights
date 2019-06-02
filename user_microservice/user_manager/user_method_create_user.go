package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/jackc/pgx"
	"golang.org/x/net/context"
)

const defaultAvatar = "default_avatar.jpg"

func (us *userManager) CreateUser(ctx context.Context, in *services.NewUserData) (*services.User, error) {

	salt, err := MakeSalt()
	if err != nil {
		logger.Errorf("Failed to make salt: %v", err.Error())
		return nil, err
	}

	password := MakePasswordHash(in.Password, salt)

	user, err := database.GetInstance().AddUser(in.Email, in.Nickname, defaultAvatar, password, salt)
	if _err, ok := err.(pgx.PgError); ok {
		logger.Errorf("Failed to add user: %v", err.Error())
		err = handlerError(_err)
		return nil, err
	}

	return &user, nil
}
