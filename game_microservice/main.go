package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/game_microservice/router"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/logger"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
	"net/http"
)

var logger = log.GetLogger("GameMS")

func init() {
	logger.SetLogLevel(logrus.Level(config.GetInt("game_ms.log_level")))
}
func main() {
	defer closer.Close()
	database.GetInstance().PopulateDatabase()

	PORT := config.GetString("game_ms.port")
	logger.Info("Game microservice started listening on", PORT)
	r := router.GetRouter()
	logger.Fatal(http.ListenAndServe(":"+PORT, r))
}
