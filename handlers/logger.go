package handlers

import (
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
)

var logger *log.Logger

func init () {
	logger = log.GetLogger("Handlers")
}
