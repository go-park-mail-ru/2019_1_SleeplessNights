package main

import (
	"flag"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/console"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"os"
	"os/exec"
)

var logger *log.Logger

func init () {
	logger = log.GetLogger("Main")
	//logger.SetLogLevel(logrus.TraceLevel)
}

func main() {
	//TODO MAKE SECRET MANAGER ?
	//TODO DELETE DATA CREATOR
	console.Title("Hello world from Sleepless Nights server!")
	console.Message("Let's check your system first...")
	//Проверяем зависимости по ПО
	softwareDependencies := []string{"psql", "consul"}
	ok := true
	for _, dep := range softwareDependencies {
		err := exec.Command(dep, "--version").Run()
		if err != nil {
			ok = false
			console.Predicate(false, dep)
		} else {
			console.Predicate(true, dep)
		}
	}
	if ok {
		console.Success("All required software is available")
	} else {
		console.Error("Some software is missing")
		return
	}

	return

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
