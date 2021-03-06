package room_manager

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/jackc/pgx"
	"google.golang.org/grpc"
)

func (rm *roomManager) CreateRoom(ctx context.Context, in *services.RoomSettings, opts ...grpc.CallOption) (*services.RoomId, error) {

	roomId, err := database.GetInstance().AddRoom(in.Talkers)
	if _err, ok := err.(pgx.PgError); ok {
		logger.Errorf("Failed to add user: %v", err.Error())
		err = handlerError(_err)
		return nil, err
	}

	r := createRoom(roomId, uint64(in.MaxConnections), in.Talkers)

	rm.Mx.Lock()
	rm.RoomsPool[roomId] = r
	rm.Mx.Unlock()

	rId := &services.RoomId{
		Id: roomId,
	}
	return rId, err
}
