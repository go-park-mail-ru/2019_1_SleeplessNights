package main

import (
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/user_microservice/user_manager"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"google.golang.org/grpc"
	"net"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("AuthMS")
	logger.SetLogLevel(logrus.TraceLevel)
}

func main() {
	defer closer.Close()

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Fatal("Auth microservice can't listen port", err)
	}

	server := grpc.NewServer()

	server.
	services.RegisterAuthCheckerServer(server, user_manager.GetInstance())

	logger.Info("Auth microservice started listening at :8081")
	err = server.Serve(lis)
	if err != nil {
		logger.Error("Auth microservice dropped with error")
		logger.Info("Restarting user_manager microservice...")
		logger.Info("Auth microservice started listening at :8081")
		err = server.Serve(lis)
	}
}
