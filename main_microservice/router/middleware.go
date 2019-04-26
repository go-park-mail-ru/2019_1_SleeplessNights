package router

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/models"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/services"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/http"
)

const (
	DomainsCORS     = "https://sleepless-nights--frontend.herokuapp.com"
	MethodsCORS     = "GET, POST, PATCH, DELETE, OPTIONS"
	CredentialsCORS = "true"
	//TODO FIX CORS HEADERS
	HeadersCORS     = "X-Requested-With, Content-type, User-Agent, Cache-Control, Cookie, Origin, Accept-Encoding, Connection, Host, Upgrade-Insecure-Requests, User-Agent, Referer, Access-Control-Request-Method, Access-Control-Request-Headers"
)

var logger *log.Logger

func init () {
	logger = log.GetLogger("Middleware")
}

type AuthHandler func(user models.User, w http.ResponseWriter, r *http.Request)

func MiddlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", DomainsCORS)
		w.Header().Set("Access-Control-Allow-Credentials", CredentialsCORS)
		w.Header().Set("Access-Control-Allow-Methods", MethodsCORS)
		w.Header().Set("Access-Control-Allow-Headers", HeadersCORS)
		next.ServeHTTP(w, r)
	})
}

func MiddlewareBasicHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func MiddlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Have some request on", r.URL)
		next.ServeHTTP(w, r)
	})
}

func MiddlewareRescue(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				fmt.Println("Unhandled handler panic:", recovered)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func MiddlewareAuth(next AuthHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Referer", r.URL.String())
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte("{}"))
			if err != nil {
				helpers.Return500(&w, err)
				return
			}
			return
		}

		grpcConn, err := grpc.Dial(
			"127.0.0.1:8081",
			grpc.WithInsecure(),
		)
		if err != nil {
			logger.Error("Can't connect to auth microservice")
			helpers.Return500(&w, err)
			return
		}
		defer grpcConn.Close()

		authManager := services.NewAuthCheckerClient(grpcConn)

		userID, err := authManager.Check(context.Background(),
			&services.SessionToken{
				Token: sessionCookie.Value,
			})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte("{}"))
			if err != nil {
				helpers.Return500(&w, err)
				return
			}
			return
		}

		user, err := database.GetInstance().GetUserViaID(userID.ID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte("{}"))
			if err != nil {
				helpers.Return500(&w, err)
				return
			}
			return
		}

		next(user, w, r)
	})
}