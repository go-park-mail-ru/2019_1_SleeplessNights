package router

import (
	forumDomain "github.com/DragonF0rm/Technopark-DBMS-Forum/domains/forum"
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter()(router *mux.Router){
	router = mux.NewRouter()
	forum := router.PathPrefix("/forum").Subrouter()
	forum.HandleFunc("/create", forumDomain.CreateHandler).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/details", forumDomain.DetailsHandler).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/create", forumDomain.CreateThreadHandler).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/users", forumDomain.UsersHandler).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/threads", forumDomain.ThreadsHandler).Methods(http.MethodGet)

	/*post    := router.PathPrefix("/post").Subrouter()
	post.HandleFunc("/{id}/details", handlers.PostGetDetails).Methods(http.MethodGet)
	post.HandleFunc("/{id}/details", handlers.PostSetDetails).Methods(http.MethodPost)

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
	user.HandleFunc("/{nickname}/profile", handlers.UserChangeProfile).Methods(http.MethodPost)*/

	router.Use(MiddlewareBasicHeaders)
	router.Use(MiddlewareLog)
	//router.Use(MiddlewareRescue)//ДОЛЖНА БЫТЬ ПОСЛЕДНЕЙ
	return
}


