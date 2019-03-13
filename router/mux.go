package router

import (

	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"github.com/gorilla/mux"
)

func GetRouter()(router *mux.Router){
	//TODO REFACTOR WITH NO VERBS
	//TODO REORGANIZE STATIC FILES ACCESS
	//TODO ADD SUBROUTERS
	router = mux.NewRouter()
<<<<<<< HEAD
	router.HandleFunc("/api/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/auth", handlers.AuthHandler).Methods("POST")//.Headers("Referer")
	router.HandleFunc("/api/profile", handlers.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/api/profile", handlers.ProfileHandler).Methods("GET")
	router.HandleFunc("/profile", handlers.ProfileUpdateHandler).Methods("PATCH")
	router.HandleFunc("/api/leaders", handlers.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/api/leaders/", handlers.LeadersHandler).Methods("GET")
=======
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	api.HandleFunc("/auth", handlers.AuthHandler).Methods("POST")//.Headers("Referer")
	api.HandleFunc("/profile", handlers.OptionsHandler).Methods("OPTIONS")
	api.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")
	api.HandleFunc("/profile", handlers.ProfileUpdateHandler).Methods("PATCH")
	api.HandleFunc("/leaders", handlers.OptionsHandler).Methods("OPTIONS")
	api.HandleFunc("/leaders", handlers.LeadersHandler).Methods("GET")
>>>>>>> 948926289d99b4b3144a838502132099524550f8
	router.HandleFunc("/img/{path}", handlers.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/img/{path}", handlers.ImgHandler).Methods("GET") //.Queries("path")

	router.Use(MiddlewareBasicHeaders)
	router.Use(MiddlewareCORS)
	return
}


