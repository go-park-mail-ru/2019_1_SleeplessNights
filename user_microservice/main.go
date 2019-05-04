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
	logger.SetLogLevel(logrus.DebugLevel)
	defer closer.Close()

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Fatal("Auth microservice can't listen port", err)
	}

	//user, err := database.GetInstance().AddUser("test@test.com", "boob", "ghghhg.img", []byte{}, []byte{})
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(user)
	//
	//user, err = database.GetInstance().AddUser("test@test.com", "boob", "ghghhg.img", []byte{}, []byte{})
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(user)

	//err = database.GetInstance().CleanerDBForTests()
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//
	//user, err = database.GetInstance().AddUser("test@test.com", "boob", "ghghhg.img", []byte{}, []byte{})
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(user)
	//
	//user.Nickname = "sdfsdf"
	//
	//err = database.GetInstance().UpdateUser(&user)
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(user)
	//
	//user.Id = 100
	//
	//err = database.GetInstance().UpdateUser(&user)
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(user)
	//
	//profile, err := database.GetInstance().GetProfile(1)
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(profile)
	//
	//profile, err = database.GetInstance().GetProfile(100)
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(profile)
	//
	//user, err = database.GetInstance().GetUserByEmail("test@test.com")
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(user)
	//
	//user, err = database.GetInstance().GetUserByEmail("test3@test.com")
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(user)
	//
	//id, password, salt, err := database.GetInstance().GetUserSignature("test@test.com")
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(id)
	//logger.Info(password)
	//logger.Info(salt)
	//
	//
	//id, password, salt, err = database.GetInstance().GetUserSignature("test3@test.com")
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//
	//user, err = database.GetInstance().AddUser("test2@test.com", "boob", "ghghhg.img", []byte{}, []byte{})
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(user)
	//
	//user, err = database.GetInstance().AddUser("test3@test.com", "boob", "ghghhg.img", []byte{}, []byte{})
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(user)
	//
	//var page services.PageData
	//page.Limit = 100
	//page.Since = 0
	//users, err := database.GetInstance().GetUsers(&page)
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(users)
	//
	//user, err = database.GetInstance().GetUserByID(1)
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(user)
	//
	//user, err = database.GetInstance().GetUserByID(100)
	//if err != nil {
	//	logger.Error(err.Error())
	//}
	//logger.Info(user)

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
