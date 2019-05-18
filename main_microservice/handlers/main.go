package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"google.golang.org/grpc"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Handlers")
	logger.SetLogLevel(logrus.Level(config.GetInt("main_ms.pkg.handlers.log_level")))
}

var userManager services.UserMSClient

func init() {
	logger.SetLogLevel(logrus.DebugLevel)
	grpcConn, err := grpc.Dial(
		"127.0.0.1:8081",
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
