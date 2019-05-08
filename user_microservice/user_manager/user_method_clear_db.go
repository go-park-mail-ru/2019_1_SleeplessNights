package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"golang.org/x/net/context"
)

func (um *userManager)ClearDB(ctx context.Context, in *services.Nothing)(*services.Nothing, error){
	//TODO RAISE ERROR UNLESS IT IS TEST CONFIGURATION
}
