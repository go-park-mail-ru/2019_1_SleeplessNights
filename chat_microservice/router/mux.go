package router

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/router/handlers"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("ChatMS")
	logger.SetLogLevel(logrus.TraceLevel)
}

func GetRouter() (router *mux.Router) {
	router = mux.NewRouter()
	ws := router.PathPrefix("/chat").Subrouter()
	ws.Handle("/connect", middleware.MiddlewareAuth(handlers.EnterChat, false))
	return
}
