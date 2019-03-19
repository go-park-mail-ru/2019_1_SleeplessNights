package router

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter()(router *mux.Router){
	router = mux.NewRouter()
	forum   := router.PathPrefix("/forum").Subrouter()
	forum.HandleFunc("/create", handlers.ForumCreate).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/create", handlers.ForumSlugCreate).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/details", handlers.ForumSlugDetails).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/threads", handlers.ForumSlugThreads).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/users", handlers.ForumSlugUsers).Methods(http.MethodGet)

	post    := router.PathPrefix("/post").Subrouter()
	post.HandleFunc("/{id}/details", handlers.PostGetDeatails).Methods(http.MethodGet)
	post.HandleFunc("/{id}/details", handlers.PostSetDeatails).Methods(http.MethodPost)

	service := router.PathPrefix("/service").Subrouter()
	service.HandleFunc("/clear", handlers.ServiceClear).Methods(http.MethodPost)
	service.HandleFunc("/status", handlers.ServiceStatus).Methods(http.MethodGet)

	thread  := router.PathPrefix("/thread").Subrouter()
	thread.HandleFunc("/{slug_or_id}/create", handlers.ThreadCreate).Methods(http.MethodPost)
	thread.HandleFunc("/{slug_or_id}/details", handlers.ThreadGetDetails).Methods(http.MethodGet)
	thread.HandleFunc("/{slug_or_id}/details", handlers.ThreadSetDetails).Methods(http.MethodPost)
	thread.HandleFunc("/{slug_or_id}/posts", handlers.ThreadPosts).Methods(http.MethodGet)
	thread.HandleFunc("/{slug_or_id}/vote", handlers.ThreadVote).Methods(http.MethodPost)

	user    := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("/{nickname}/create", handlers.UserCreate).Methods(http.MethodPost)
	user.HandleFunc("/{nickname}/profile", handlers.UserGetProfile).Methods(http.MethodGet)
	user.HandleFunc("/{nickname}/profile", handlers.UserChangeProfile).Methods(http.MethodPost)

	router.Use(MiddlewareBasicHeaders)
	router.Use(MiddlewareRescue)
	router.Use(MiddlewareLog)
	return
}


