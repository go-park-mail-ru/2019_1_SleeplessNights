package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/router"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
	"github.com/xlab/closer"
	"net/http"
	"os"
	"sync"
)

var logger *log.Logger

func init () {
	logger = log.GetLogger("Main")
	//logger.SetLogLevel(logrus.TraceLevel)
}

func main () {
	defer closer.Close()
	//faker.CreateFakeData(10)
	//faker.CreateFakePacks()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	logger.Info("Main microservice started listening on", PORT)
	r := router.GetRouter()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		logger.Fatal(http.ListenAndServe(":"+PORT, r))
		wg.Done()
	}(&wg)

	/*user, _ := database.GetInstance().GetUserViaID(1)
	cookie, _ := auth_microservice.MakeSession(user)
	connUser := exec.Command(`./scripts/ws-connect.sh`, PORT, cookie.Value)
	err := connUser.Run()
	if err != nil {
		logger.Error(err)
	}
	user, _ = database.GetInstance().GetUserViaID(2)
	cookie, _ = auth_microservice.MakeSession(user)
	connUser = exec.Command(`./scripts/ws-connect.sh`, PORT, cookie.Value)
	err = connUser.Run()
	if err != nil {
		logger.Error(err)
	}*/

	wg.Wait()
}