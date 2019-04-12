package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/faker"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"github.com/xlab/closer"
	"net/http"
	"os"
)

func main() {
	//TODO MAKE SECRET MANAGER ?

	//TODO DELETE DATA CREATOR

	defer closer.Close()

	logger.Info.Printf("\nSuccessfully connected to database on: %s", "...") //TODO PORT

	faker.CreateFakeData(10)
	faker.CreateFakePacks()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	logger.Info.Println("Started listening on", PORT)
	r := router.GetRouter()
	logger.Fatal.Println(http.ListenAndServe(":"+PORT, r))
}
