package router

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter()(router *mux.Router){
	//TODO REFACTOR WITH NO VERBS
	//TODO REORGANIZE STATIC FILES ACCESS
	//TODO ADD AMAZON S3
	//TODO ADD RECOVER MIDDLEWARE
	router = mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", handlers.RegisterHandler).Methods(http.MethodPost)
	api.HandleFunc("/auth", handlers.AuthHandler).Methods(http.MethodPost)
	api.HandleFunc("/profile", handlers.OptionsHandler).Methods(http.MethodOptions)
	api.HandleFunc("/profile", handlers.ProfileHandler).Methods(http.MethodGet)
	api.HandleFunc("/profile", handlers.ProfileUpdateHandler).Methods(http.MethodPatch)
	api.HandleFunc("/leaders", handlers.OptionsHandler).Methods(http.MethodOptions)
	api.HandleFunc("/leaders", handlers.LeadersHandler).Methods(http.MethodGet)
	router.HandleFunc("/img/{path}", handlers.OptionsHandler).Methods(http.MethodOptions)
	router.HandleFunc("/img/{path}", handlers.ImgHandler).Methods(http.MethodGet)

	router.Use(MiddlewareBasicHeaders)
	router.Use(MiddlewareCORS)
	router.Use(MiddlewareLog)
	return
}

git
