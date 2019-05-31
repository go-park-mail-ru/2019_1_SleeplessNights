package router

import (
	"github.com/gorilla/mux"
	"reflect"
	"testing"
)

func TestGetRouter(t *testing.T) {
	if reflect.TypeOf(GetRouter()) != reflect.PtrTo(reflect.TypeOf(mux.Router{})) {
		t.Errorf("GetInstance method returns value with wrong type: got %s, whant %s",
			reflect.TypeOf(reflect.TypeOf(GetRouter())),
			reflect.PtrTo(reflect.TypeOf(mux.Router{})))
	}
}

