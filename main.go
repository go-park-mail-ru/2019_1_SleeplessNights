package main

import (
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
)

var logger *log.Logger

func init () {
	logger = log.GetLogger("Main")
	//logger.SetLogLevel(logrus.TraceLevel)
}

func main() {
	//TODO MAKE SECRET MANAGER ?
	//TODO DELETE DATA CREATOR

}
