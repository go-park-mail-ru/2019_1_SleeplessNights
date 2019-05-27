package user_manager

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/errors"
	"github.com/jackc/pgx"
	"reflect"
	"testing"
)

func TestGetInstance(t *testing.T) {
	if reflect.TypeOf(GetInstance()) != reflect.PtrTo(reflect.TypeOf(userManager{})) {
		t.Errorf("GetInstance method returns value with wrong type: got %s, whant %s",
			reflect.TypeOf(reflect.TypeOf(GetInstance())),
			reflect.PtrTo(reflect.TypeOf(userManager{})))
	}
}

func TestHandlerErrorSuccessful(t *testing.T) {
	var err pgx.PgError
	err.Code = uniqueViolation
	if _err := handlerError(err); _err != errors.DataBaseUniqueViolation {
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