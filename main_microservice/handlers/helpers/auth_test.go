package helpers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"testing"
)

func TestBuildSessionCookie(t *testing.T) {
	session := &services.SessionToken{}
	_ = BuildSessionCookie(session)
}
