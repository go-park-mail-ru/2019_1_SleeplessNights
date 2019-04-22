package main

import (
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
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

func main() {
	//TODO MAKE SECRET MANAGER ?
	//TODO DELETE DATA CREATOR

	/*logger.Trace("Hello world, it's func main() from main.go")
	logger.Debug("If I don't work correctly, you can use this log level to debug me")
	logger.Info("Through this I will send you feedback. For example PI =", 3.14)
	logger.Warningf("O-oh, there are only %d numbers after decimal separator", 2)
	errorFields := make(map[string]interface{}, 1)
	errorFields["reason"] = "PI number accuracy"
	logger.ErrorWithFields(errorFields, "ERROR: we need more of that")
	logger.Fatal("Can not continue")*/

	defer closer.Close()
	//faker.CreateFakeData(10)
	//faker.CreateFakePacks()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	logger.Info("Started listening on", PORT)
	r := router.GetRouter()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		logger.Fatal(http.ListenAndServe(":"+PORT, r))
		wg.Done()
	}(&wg)

	/*user, _ := database.GetInstance().GetUserViaID(1)
	cookie, _ := auth.MakeSession(user)
	connUser := exec.Command(`./ws-connect.sh`, PORT, cookie.Value)
	err := connUser.Run()
	if err != nil {
		logger.Error(err)
	}
	user, _ = database.GetInstance().GetUserViaID(2)
	cookie, _ = auth.MakeSession(user)
	connUser = exec.Command(`./ws-connect.sh`, PORT, cookie.Value)
	err = connUser.Run()
	if err != nil {
		logger.Error(err)
	}*/

	wg.Wait()
}
