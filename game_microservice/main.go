package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database/faker"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/router"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"net/http"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("ChatMS")
	logger.SetLogLevel(logrus.TraceLevel)
}

func main() {
	defer closer.Close()

	/*fmt.Println("Achtung! Questions are about to be deleted. Proceed?")
	fmt.Println("Y/N")
	flag := ""
	_, err := fmt.Scanln(flag)
	if err != nil {
		return
	}
	if (flag != "Y") && (flag != "y") {
		return
	}*/
	err := database.GetInstance().CleanerDBForTests()
	if err != nil {
		logger.Errorf(err.Error())
	}

	faker.CreateFakePacks()

	PORT := "8006"
	logger.Info("Game microservice started listening on", PORT)
	r := router.GetRouter()
	logger.Fatal(http.ListenAndServe(":"+PORT, r))
}
