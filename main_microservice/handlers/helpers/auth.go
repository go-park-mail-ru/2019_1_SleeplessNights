package helpers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"net/http"
	"time"
)

const defaultCookieLifeLen = 4 * time.Hour

func BuildSessionCookie(sessionToken *services.SessionToken)(sessionCookie http.Cookie) {
	cookieLifeLen, err := time.ParseDuration(config.GetString("main_ms.pkg.helpers.cookie.life_len"))
	if err != nil {
		cookieLifeLen = defaultCookieLifeLen
	}
	sessionCookie = http.Cookie{
		Name: config.GetString("main_ms.pkg.helpers.cookie.name"),
		Value: sessionToken.Token,
		Expires: time.Now().Add(cookieLifeLen),
		HttpOnly: config.GetBool("main_ms.pkg.helpers.cookie.http_only"),
	Secure: config.GetBool("main_ms.pkg.helpers.cookie.secure"),
	}
	return
}
