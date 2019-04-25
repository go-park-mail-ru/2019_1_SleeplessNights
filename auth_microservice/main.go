package auth_microservice

import (
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/services"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Auth")
	logger.SetLogLevel(logrus.TraceLevel)
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Fatal("Auth microservice can't listen port", err)
	}

	server := grpc.NewServer()
	services.RegisterAuthCheckerServer(server, GetInstance())

	logger.Info("Auth microservice started listening at :8081")
	err = server.Serve(lis)
	if err != nil {
		logger.Error("Auth microservice dropped with error")
		logger.Info("Restarting auth microservice...")
		logger.Info("Auth microservice started listening at :8081")
		err = server.Serve(lis)
	}
}
