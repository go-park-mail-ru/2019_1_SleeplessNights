package main

import (
	"flag"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"os"
)

var logger *log.Logger

func init () {
	logger = log.GetLogger("Main")
	//logger.SetLogLevel(logrus.TraceLevel)
}

func main() {
	//TODO MAKE SECRET MANAGER ?
	//TODO DELETE DATA CREATOR

	conf := flag.String("-conf","DEV", "Sets the configuration to: DEV (default), TEST, LOCAL or PROD")
	verb := flag.Bool("v", false, "Shows more info during the execution")
	flag.Parse()

	configs := make(map[string]struct{})
	configs["DEV"]   = struct{}{}
	configs["TEST"]  = struct{}{}
	configs["LOCAL"] = struct{}{}
	configs["PROD"]  = struct{}{}
	if _, found := configs[*conf]; !found {
		logger.Fatal("Unexpected configuration key", *conf, "provided")
	}
	err := os.Setenv("CONFIG", *conf)
	if err != nil {
		logger.Fatal("Failed to set configuration env variable", err)
	}

	if *verb {
		logger.Debug("Setting config params to more verbose")
	}

	/*err = exec.Command("go run auth_microservice/main.go").Run()
	if err != nil {
		logger.Fatal("Can't run auth microservice:", err)
	}*/
	/*err = exec.Command("go","run", "main_microservice/main.go").Run()
	if err != nil {
		logger.Fatal("Can't run main microservice")
	}*/
}
