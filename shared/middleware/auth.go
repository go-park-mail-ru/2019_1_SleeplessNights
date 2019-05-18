package middleware

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/xlab/closer"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/http"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Middleware")
}

var userManager services.UserMSClient

func init() {
	var err error
	grpcConn, err := grpc.Dial(
		config.GetString("user_ms.address"),
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Fatal("Can't connect to auth microservice via grpc")
	}
	userManager = services.NewUserMSClient(grpcConn)
	closer.Bind(func() {
		err := grpcConn.Close()
		if err != nil {
			logger.Error("Error occurred while closing grpc connection", err)
		}
	})
}

type AuthHandler func(user *services.User, w http.ResponseWriter, r *http.Request)

func MiddlewareAuth(next AuthHandler, strict bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Referer", r.URL.String())
		sessionCookie, err := r.Cookie("session_token")
		if err != nil {
			logger.Info("Got connection with missing cookie")
			if err == http.ErrNoCookie && !strict {
				logger.Info("Trying to connect anonymous client")

				user := services.User{
					Id:         0,
					Email:      "",
					Nickname:   "Guest",
					AvatarPath: "default_avatar.jpg",
				}
				next(&user, w, r)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte("{}"))
			if err != nil {
				//helpers.Return500(&w, err)
				(w).WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		user, err := userManager.CheckToken(context.Background(),
			&services.SessionToken{
				Token: sessionCookie.Value,
			})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err = w.Write([]byte("{}"))
			if err != nil {
				//helpers.Return500(&w, err)
				(w).WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}

		next(user, w, r)
	})
}
