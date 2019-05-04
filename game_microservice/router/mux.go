package router

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/router/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"
	"github.com/gorilla/mux"
)

func GetRouter() (router *mux.Router) {

	router = mux.NewRouter()
	ws := router.PathPrefix("/ws").Subrouter()
	ws.Handle("/connect", middleware.MiddlewareAuth(handlers.UpgradeWs, true))

	router.Use(middleware.MiddlewareBasicHeaders)
	router.Use(middleware.MiddlewareCORS)
	router.Use(middleware.MiddlewareLog)
	router.Use(middleware.MiddlewareRescue)
	return
}
