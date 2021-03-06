package router

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter() (router *mux.Router) {
	//TODO REORGANIZE STATIC FILES ACCESS
	//TODO ADD AMAZON S3
	router = mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/user", handlers.RegisterHandler).Methods(http.MethodPost)
	api.HandleFunc("/user", handlers.OptionsHandler).Methods(http.MethodOptions)
	api.HandleFunc("/session", handlers.AuthHandler).Methods(http.MethodPost)
	api.HandleFunc("/session", handlers.AuthDeleteHandler).Methods(http.MethodDelete)
	api.HandleFunc("/session", handlers.OptionsHandler).Methods(http.MethodOptions)
	api.HandleFunc("/leader", handlers.OptionsHandler).Methods(http.MethodOptions)
	api.HandleFunc("/leader", handlers.LeadersHandler).Methods(http.MethodGet)
	api.HandleFunc("/profile", handlers.OptionsHandler).Methods(http.MethodOptions)
	//Запросы, требующие авторизации
	api.Handle("/profile", middleware.MiddlewareAuth(handlers.ProfileHandler, true)).Methods(http.MethodGet)
	api.Handle("/profile", middleware.MiddlewareAuth(handlers.ProfileUpdateHandler, true)).Methods(http.MethodPatch)

	router.HandleFunc("/img/{path}", handlers.OptionsHandler).Methods(http.MethodOptions)
	router.HandleFunc("/img/{path}", handlers.ImgHandler).Methods(http.MethodGet)

	router.Use(middleware.MiddlewareBasicHeaders)
	router.Use(middleware.MiddlewareCORS)
	router.Use(middleware.MiddlewareLog)
	router.Use(middleware.MiddlewareRescue)
	return
}
