package router

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/handlers"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Router")
	logger.SetLogLevel(logrus.TraceLevel)
}

func GetRouter() (router *mux.Router) {
	router = mux.NewRouter()
	ws := router.PathPrefix("/chat").Subrouter()

	ws.Handle("/connect", middleware.MiddlewareAuth(handlers.EnterChat, false))

	router.Use(middleware.MiddlewareBasicHeaders)
	router.Use(middleware.MiddlewareCORS)
	router.Use(middleware.MiddlewareLog)
	router.Use(middleware.MiddlewareRescue)
	return
}
