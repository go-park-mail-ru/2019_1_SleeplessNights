package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database/models"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
)

var logger *log.Logger

func init() {
	logger = log.GetLogger("ChatMS")
	logger.SetLogLevel(logrus.TraceLevel)
}

func main() {
	defer closer.Close()

	err := database.GetInstance().AddQuestionPack("something")
	if err != nil {
		logger.Fatal(err.Error())
	}

	q := models.Question{
		Answers: []string{"sdfsdf", "sdfdsfdsfdsf", "hjkhjkhjkh", "werewrhytu"},
		Correct: 1,
		Text: "vnkuerlNCSDNVKSVNSDV",
		PackID: 0,
	}

	err = database.GetInstance().AddQuestion(q)
	if err != nil {
		logger.Fatal(err.Error())
	}

	//In case of a lack of data, break parentheses
	//err := database.GetInstance().CleanerDBForTests()
	//if err != nil {
	//	logger.Errorf(err.Error())
	//}

	//faker.CreateFakePacks()

	//PORT := "8006"
	//logger.Info("Game microservice started listening on", PORT)
	//r := router.GetRouter()
	//logger.Fatal(http.ListenAndServe(":"+PORT, r))
}
