package helpers_test

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"testing"
)

func TestReturn400(t *testing.T) {
	var requestErrors helpers.ErrorSet
	_, err := helpers.MarshalToJSON(requestErrors)
	if err != nil {
		t.Errorf("Errors returned error: %v\n", err.Error())
	}
}
