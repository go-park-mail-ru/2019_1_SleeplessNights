package helpers

import (
	"context"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/services"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

func BuildSessionCookie(userID uint64)(sessionCookie http.Cookie, err error) {
	grpcConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Error("Can't connect to auth microservice")
		return
	}
	defer grpcConn.Close()

	authManager := services.NewAuthCheckerClient(grpcConn)

	sessionToken, err := authManager.MakeToken(context.Background(),
		&services.UserID{
			ID: userID,
		})
	if err != nil {
		logger.Error("Can't make token for this user")
		return
	}

	sessionCookie = http.Cookie{
		Name: "session_token",
		Value: sessionToken.Token,
		Expires: time.Now().Add(4 * time.Hour),
		HttpOnly: true,
	}
	return
}
