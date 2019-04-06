package router

import (
	forumDomain  "github.com/DragonF0rm/Technopark-DBMS-Forum/domains/forum"
	postDomain   "github.com/DragonF0rm/Technopark-DBMS-Forum/domains/post"
	statusDomain "github.com/DragonF0rm/Technopark-DBMS-Forum/domains/service"
	threadDomain "github.com/DragonF0rm/Technopark-DBMS-Forum/domains/thread"
	userDomain   "github.com/DragonF0rm/Technopark-DBMS-Forum/domains/user"
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter()(router *mux.Router){
	router = mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	forum := api.PathPrefix("/forum").Subrouter()
	forum.HandleFunc("/create", forumDomain.CreateHandler).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/details", forumDomain.DetailsHandler).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/create", forumDomain.CreateThreadHandler).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/users", forumDomain.UsersHandler).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/threads", forumDomain.ThreadsHandler).Methods(http.MethodGet)

	post := api.PathPrefix("/post").Subrouter()
	post.HandleFunc("/{id}/details", postDomain.DetailsHandler).Methods(http.MethodGet)
	post.HandleFunc("/{id}/details", postDomain.EditHandler).Methods(http.MethodPost)

	service := api.PathPrefix("/service").Subrouter()
	service.HandleFunc("/clear", statusDomain.ClearHandler).Methods(http.MethodPost)
	service.HandleFunc("/status", statusDomain.StatusHandler).Methods(http.MethodGet)

	thread  := api.PathPrefix("/thread").Subrouter()
	thread.HandleFunc("/{slug_or_id}/create", threadDomain.CreatePostsHandler).Methods(http.MethodPost)
	thread.HandleFunc("/{slug_or_id}/details", threadDomain.DetailsHandler).Methods(http.MethodGet)
	thread.HandleFunc("/{slug_or_id}/details", threadDomain.UpdateHandler).Methods(http.MethodPost)
	thread.HandleFunc("/{slug_or_id}/posts", threadDomain.PostsHandler).Methods(http.MethodGet)
	thread.HandleFunc("/{slug_or_id}/vote", threadDomain.VoteHandler).Methods(http.MethodPost)

	user := api.PathPrefix("/user").Subrouter()
	user.HandleFunc("/{nickname}/create", userDomain.CreateHandler).Methods(http.MethodPost)
	user.HandleFunc("/{nickname}/profile", userDomain.ProfileHandler).Methods(http.MethodGet)
	user.HandleFunc("/{nickname}/profile", userDomain.EditProfileHandler).Methods(http.MethodPost)

	router.Use(MiddlewareBasicHeaders)
	router.Use(MiddlewareLog)
	//router.Use(MiddlewareRescue)//ДОЛЖНА БЫТЬ ПОСЛЕДНЕЙ
	return
}


