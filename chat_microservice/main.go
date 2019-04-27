
package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/chat_microservice/router"
	log "github.com/go-park-mail-ru/2019_1_SleeplessNights/meta/logger"
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

	id, err := database.GetInstance().UpdateUser(1, "poll", "/home")
	if err != nil{
		fmt.Print(err.Error())
	}
	fmt.Print(id)

	id, err = database.GetInstance().CreateRoom([]uint64{2,3})
	if err != nil{
		fmt.Print(err.Error())
	}
	fmt.Print(id)

	PORT := "8005"

	logger.Info("Chat microservice started listening on", PORT)
	r := router.GetRouter()

	logger.Fatal(http.ListenAndServe(":"+PORT, r))

}

