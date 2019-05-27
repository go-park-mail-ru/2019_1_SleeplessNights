package database

import (
	"reflect"
	"testing"
)

func TestGetInstance(t *testing.T) {
	if reflect.TypeOf(GetInstance()) != reflect.PtrTo(reflect.TypeOf(dbManager{})) {
		t.Errorf("GetInstance method returns value with wrong type: got %s, whant %s",
			reflect.TypeOf(reflect.TypeOf(GetInstance())),
			reflect.PtrTo(reflect.TypeOf(dbManager{})))
	}
}