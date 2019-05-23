package room_manager

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/jackc/pgx"
	"google.golang.org/grpc"
)

func (rm *roomManager) DeleteRoom(ctx context.Context, in *services.Room, opts ...grpc.CallOption) (n *services.Nothing, err error) {

	if in.Id == globalChatId {
		logger.Errorf(`Failed to delete global chat`)
		err = errors.New("ERROR: trying delete global chat")
		return
	}

	err = database.GetInstance().DeleteRoom(in.Id)
	if _err, ok := err.(pgx.PgError); ok {
		logger.Errorf("Failed to add user: %v", err.Error())
		err = handlerError(_err)
		return nil, err
	}

	rm.Mx.Lock()
	delete(chat.RoomsPool, in.Id)
	rm.Mx.Unlock()

	return
}
