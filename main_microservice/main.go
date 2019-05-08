package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/router"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/xlab/closer"
	"net/http"
	"os"
	"sync"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("Main")
	//logger.SetLogLevel(logrus.TraceLevel)
}

func main() {
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

	/*user_manager, _ := database.GetInstance().GetUserByID(1)
	cookie, _ := user_microservice.MakeSession(user_manager)
	connUser := exec.Command(`./scripts/ws-connect.sh`, PORT, cookie.Value)
	err := connUser.Run()
	if err != nil {
		logger.Error(err)
	}
	user_manager, _ = database.GetInstance().GetUserByID(2)
	cookie, _ = user_microservice.MakeSession(user_manager)
	connUser = exec.Command(`./scripts/ws-connect.sh`, PORT, cookie.Value)
	err = connUser.Run()
	if err != nil {
		logger.Error(err)
	}*/

	wg.Wait()
}
