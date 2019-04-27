package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/router"

	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"net/http"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("ChatMS")
	logger.SetLogLevel(logrus.TraceLevel)
}

func main() {
	defer closer.Close()

	PORT := "8005"

	logger.Info("Chat microservice started listening on", PORT)
	r := router.GetRouter()

	logger.Fatal(http.ListenAndServe(":"+PORT, r))

}
