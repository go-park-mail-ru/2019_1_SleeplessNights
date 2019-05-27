package room_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	"github.com/jackc/pgx"
	"reflect"
	"testing"
)

func TestGetInstance(t *testing.T) {
	if reflect.TypeOf(GetInstance()) != reflect.PtrTo(reflect.TypeOf(roomManager{})) {
		t.Errorf("GetInstance method returns value with wrong type: got %s, whant %s",
			reflect.TypeOf(reflect.TypeOf(GetInstance())),
			reflect.PtrTo(reflect.TypeOf(roomManager{})))
	}
}

func TestCreateRoomSuccessful(t *testing.T) {
	id := uint64(1)
	maxConn := uint64(2)
	talkersArray := []uint64{1, 2}

	r := createRoom(id, maxConn, talkersArray)
	if r.Id != id {
		t.Errorf("CreateRoom returned wrong value:\ngot %v\nwhant %v",
			r.Id, id)
	}
}

func TestHandlerErrorSuccessful(t *testing.T) {
	var err pgx.PgError
	err.Code = foreignKeyViolation
	if _err := handlerError(err); _err != errors.DataBaseForeignKeyViolation {
		t.Errorf("HandlerError returned wrong error:\ngot %v\nwhant %v",
			_err, errors.DataBaseForeignKeyViolation)
	}

	err.Code = nodataFound
	if _err := handlerError(err); _err != errors.DataBaseNoDataFound {
		t.Errorf("HandlerError returned wrong error:\ngot %v\nwhant %v",
			_err, errors.DataBaseNoDataFound)
	}

	var someErr pgx.PgError
	if _err := handlerError(someErr); _err != someErr {
		t.Errorf("HandlerError returned wrong error:\ngot %v\nwhant %v",
			_err, errors.DataBaseNoDataFound)
	}
}
