package helpers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"net/http"
	"time"
)

func BuildSessionCookie(sessionToken *services.SessionToken)(sessionCookie http.Cookie) {
	sessionCookie = http.Cookie{
		Name: "session_token",
		Value: sessionToken.Token,
		Expires: time.Now().Add(4 * time.Hour),
		HttpOnly: true,
	}
	return
}
