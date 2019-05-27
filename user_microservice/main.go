package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
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
	logger = log.GetLogger("UM_AuthMS")
	logger.SetLogLevel(logrus.Level(config.GetInt("user_ms.log_level")))
}

func main() {
	logger.SetLogLevel(logrus.DebugLevel)
	defer closer.Close()

	lis, err := net.Listen("tcp", config.GetString("user_ms.address"))
	if err != nil {
		logger.Fatal("Auth microservice can't listen port", err)
	}

	server := grpc.NewServer()

	services.RegisterUserMSServer(server, user_manager.GetInstance())

	logger.Info("Auth microservice started listening at :8081")
	err = server.Serve(lis)
	if err != nil {
		logger.Error("Auth microservice dropped with error")
		logger.Info("Restarting user_manager microservice...")
		logger.Info("Auth microservice started listening at :8081")
		err = server.Serve(lis)
	}
}
