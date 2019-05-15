package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/database"
	"github.com/jackc/pgx"
	"golang.org/x/net/context"
)

func (us *userManager)ClearDB(ctx context.Context, in *services.Nothing)(*services.Nothing, error){
	//TODO RAISE ERROR UNLESS IT IS TEST CONFIGURATION
	err := database.GetInstance().CleanerDBForTests()
	_err, _ := err.(pgx.PgError)
	err = handlerError(_err)
	return in, err
}
