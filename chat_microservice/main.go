package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/router"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"net/http"
)
//go:generate echo "Hello world"
var logger *log.Logger

func init() {
	logger = log.GetLogger("ChatMS")
	logger.SetLogLevel(logrus.Level(config.GetInt("chat_ms.log_level")))
}

func main() {
	defer closer.Close()

	PORT := config.GetString("chat_ms.port")
	logger.Info("Chat microservice started listening on", PORT)
	r := router.GetRouter()

	logger.Fatal(http.ListenAndServe(":"+PORT, r))

}
