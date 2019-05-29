package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/router"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"net/http"
	"sync"
)
//go:generate $GOPATH/bin/easyjson ./...
var logger *log.Logger

func init() {
	logger = log.GetLogger("Main")
	logger.SetLogLevel(logrus.Level(config.GetInt("main_ms.log_level")))
}

func main() {
	defer closer.Close()
	PORT := config.GetString("main_ms.port")
	r := router.GetRouter()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		logger.Fatal(http.ListenAndServe(":"+string(PORT), r))
		wg.Done()
	}(&wg)
	logger.Info("Main microservice started listening on", PORT)
	//Здесь можно вызвать скрипты
	wg.Wait()
}
