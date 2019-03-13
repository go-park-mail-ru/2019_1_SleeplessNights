package router

import (

	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	DomainsCORS     = "https://sleeples-nights--frontend.herokuapp.com"
	MethodsCORS     = "GET POST PATCH OPTIONS"
	CredentialsCORS = "true"
	//TODO FIX CORS HEADERS
	HeadersCORS     = "X-Requested-With, Content-type, User-Agent, Cache-Control, Cookie, Origin, Accept-Encoding, Connection, Host, Upgrade-Insecure-Requests, User-Agent, Referer, Access-Control-Request-Method, Access-Control-Request-Headers"
)


func GetRouter()(router *mux.Router){
	//TODO REFACTOR WITH NO VERBS
	//TODO ADD SUBROUTERS
	router = mux.NewRouter()
	router.HandleFunc("/api/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/auth", handlers.AuthHandler).Methods("POST")//.Headers("Referer")
	router.HandleFunc("/api/profile", handlers.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/api/profile", handlers.ProfileHandler).Methods("GET")
	router.HandleFunc("/profile", handlers.ProfileUpdateHandler).Methods("PATCH")
	router.HandleFunc("/api/leaders", handlers.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/api/leaders/{page}", handlers.LeadersHandler).Methods("GET")
	router.HandleFunc("/img/{path}", handlers.OptionsHandler).Methods("OPTIONS")
	router.HandleFunc("/img/{path}", handlers.ImgHandler).Methods("GET") //.Queries("path")
	return
}

func SetBasicHeaders(w *http.ResponseWriter){
	//TODO REWRITE AS MIDDLEWARE
	(*w).Header().Set("Access-Control-Allow-Origin", DomainsCORS)
	(*w).Header().Set("Access-Control-Allow-Credentials", CredentialsCORS)
	(*w).Header().Set("Access-Control-Allow-Methods", MethodsCORS)
	(*w).Header().Set("Access-Control-Allow-Headers", HeadersCORS)
	(*w).Header().Set("Content-type", "application/json")
}