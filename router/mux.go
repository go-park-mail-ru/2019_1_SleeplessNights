package router

import (

	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"github.com/gorilla/mux"
)

func GetRouter()(router *mux.Router){
	//TODO REFACTOR WITH NO VERBS
	//TODO REORGANIZE STATIC FILES ACCESS
	//TODO ADD AMAZON S3
	//TODO ADD RECOVER MIDDLEWARE
	router = mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	api.HandleFunc("/auth", handlers.AuthHandler).Methods("POST")
	api.HandleFunc("/profile", handlers.OptionsHandler).Methods("OPTIONS")
	api.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")
	api.HandleFunc("/profile", handlers.ProfileUpdateHandler).Methods("PATCH")
	api.HandleFunc("/leaders", handlers.OptionsHandler).Methods("OPTIONS")
	api.HandleFunc("/leaders", handlers.LeadersHandler).Methods("GET")
	router.HandleFunc("/img/{path}", handlers.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/img/{path}", handlers.ImgHandler).Methods("GET")

	router.Use(MiddlewareBasicHeaders)
	router.Use(MiddlewareCORS)
	router.Use(MiddlewareLog)
	return
}


