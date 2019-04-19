package router

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter() (router *mux.Router) {
	//TODO REWORK AUTH AND MOVE IT TO MIDDLEWARE
	//TODO REORGANIZE STATIC FILES ACCESS
	//TODO ADD AMAZON S3
	router = mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/user", handlers.RegisterHandler).Methods(http.MethodPost)
	api.HandleFunc("/session", handlers.AuthHandler).Methods(http.MethodPost)
	api.HandleFunc("/session", handlers.AuthDeleteHandler).Methods(http.MethodDelete)
	api.HandleFunc("/profile", handlers.OptionsHandler).Methods(http.MethodOptions)
	api.HandleFunc("/leader", handlers.OptionsHandler).Methods(http.MethodOptions)
	api.HandleFunc("/leader", handlers.LeadersHandler).Methods(http.MethodGet)
	//Запросы, требующие авторизации
	api.Handle("/profile", MiddlewareAuth(handlers.ProfileHandler)).Methods(http.MethodGet)
	api.Handle("/profile", MiddlewareAuth(handlers.ProfileUpdateHandler)).Methods(http.MethodPatch)

	ws := router.PathPrefix("/ws").Subrouter()
	ws.Handle("/connect", MiddlewareAuth(handlers.UpgradeWs))

	router.HandleFunc("/img/{path}", handlers.OptionsHandler).Methods(http.MethodOptions)
	router.HandleFunc("/img/{path}", handlers.ImgHandler).Methods(http.MethodGet)

	router.Use(MiddlewareBasicHeaders)
	router.Use(MiddlewareCORS)
	router.Use(MiddlewareLog)
	router.Use(MiddlewareRescue)
	return
}
