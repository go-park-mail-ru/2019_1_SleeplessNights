package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/router"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"net/http"
)

var logger = log.GetLogger("GameMS")

func init() {
	logger.SetLogLevel(logrus.TraceLevel)
}

func main() {
	defer closer.Close()

	PORT := "8006"
	logger.Info("Game microservice started listening on", PORT)
	r := router.GetRouter()
	logger.Fatal(http.ListenAndServe(":"+PORT, r))
}
