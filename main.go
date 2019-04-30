package main

import (
	"flag"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/console"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/xlab/closer"
	"strconv"
	"time"

	//"github.com/opencopilot/consulkvjson"
	//"io/ioutil"
	"os"
	"os/exec"
)

var logger *log.Logger

func init () {
	logger = log.GetLogger("Main")
	//logger.SetLogLevel(logrus.TraceLevel)
}

var jobs map[string]int

/*func startJob(script, jobName string)(err error) {
	err = exec.Command(script).Run()
	if err != nil {
		return
	}

	id, err := strconv.Atoi(os.Getenv("!"))
	if err != nil {
		return
	}

	output, err := exec.Command(os.Getenv(microservices.BASEPATH) +"/scripts/get-job-id.sh").Output()
	if err != nil {
		console.Error("Failed to get id of job started", script)
		return err
	}

	id, err := strconv.Atoi(
		strings.TrimPrefix(
			strings.TrimSuffix(
				strings.Split(string(output), " ")[0],
				"]"),
				"["))


	jobs[jobName]=id
	return
}*/

func main() {
	//TODO MAKE SECRET MANAGER ?
	//TODO DELETE DATA CREATOR
	defer console.Title("Good bye!")
	defer closer.Close()
	/*closer.Bind(func() {
		console.Title("Terminating child jobs...")
		ok := true
		for jobName, job := range jobs {
			err := exec.Command("pkill", "-P", strconv.Itoa(job))
			console.Predicate(err != nil,"%s %d", jobName, job)
			if err != nil {
				ok = false
			}
		}
		if ok {
			console.Success("All child jobs terminated")
		} else {
			console.Error("Child jobs termination completed wit errors")
		}
	})*/

	console.Title("Hello world from Sleepless Nights server!")
	console.Message("Let's check your system first...")
	//Проверяем зависимости по ПО

	softwareDependencies := make(map[string]string)
	softwareDependencies["psql"]   = "--version"
	softwareDependencies["consul"] = "--version"
	softwareDependencies["xterm"]  = "-help"
	softwareDependencies["jq"]     = "--version"
	ok := true
	for dep, arg := range softwareDependencies {
		err := exec.Command(dep, arg).Run()
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

	err := os.Setenv(config.BASEPATH, os.Getenv("PWD"))
	if err != nil {
		console.Error("Failed to set BASEPATH")
		return
	}
	console.Message("BASEPATH set to %s", os.Getenv(config.BASEPATH))

	console.Title("Starting consul...")
	consulCmd := exec.Command(os.Getenv(config.BASEPATH) + "/scripts/consul-start.sh")
	consulStdin, err := consulCmd.StdinPipe()
	if err != nil {
		logger.Fatal("Failed to get consul stdin pipe")
	}
	err = consulCmd.Start()
	if err != nil {
		console.Error("Failed to start consul")
		logger.Fatal(err)
	}
	closer.Bind(func() {
		err := consulStdin.Close()
		if err != nil {
			logger.Error("Consul finished with error:",err.Error())
		}
		console.Message("Consul finished")
	})
	console.Success("Consul started")

	for {
		counter := 0
		time.Sleep(time.Second)
		console.Message(strconv.Itoa(counter))
		counter++
	}
	return

	//jsonConfig, err := ioutil.ReadFile("microservices.json")
	if err != nil {
		console.Error("Failed to read microservices.json")
		return
	}

	//kvs, err := consulkvjson.ToKVs(jsonConfig)
	if err != nil {
		console.Error("Failed to convert JSON microservices into consul KV's")
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
	err = os.Setenv("CONFIG", *conf)
	if err != nil {
		logger.Fatal("Failed to set configuration env variable", err)
	}

	if *verb {
		logger.Debug("Setting microservices params to more verbose")
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
