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



	/*err := exec.Command("go","run", "auth_microservice/main.go").Run()
	if err != nil {
		logger.Fatal("Can't run auth microservice:", err)
	}
	err = exec.Command("go","run", "main_microservice/main.go").Run()
	if err != nil {
		logger.Fatal("Can't run main microservice")
	}*/
}
