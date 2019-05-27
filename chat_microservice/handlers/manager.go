package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Handlers")
	logger.SetLogLevel(logrus.Level(config.GetInt("chat_ms.log_level")))
}
