package router

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"
	"github.com/gorilla/mux"
)

func GetRouter() (router *mux.Router) {
	router = mux.NewRouter()
	ws := router.PathPrefix("/api").Subrouter()

	ws.Handle("/chat", middleware.MiddlewareAuth(handlers.EnterChat, false))

	router.Use(middleware.MiddlewareBasicHeaders)
	router.Use(middleware.MiddlewareCORS)
	router.Use(middleware.MiddlewareLog)
	router.Use(middleware.MiddlewareRescue)
	return
}
